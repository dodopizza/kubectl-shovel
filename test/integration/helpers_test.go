//go:build integration
// +build integration

package integration_test

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"

	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s "k8s.io/client-go/kubernetes"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/dodopizza/kubectl-shovel/internal/globals"
	"github.com/dodopizza/kubectl-shovel/internal/kubernetes"
)

var (
	sampleAppName       = "sample-app"
	sampleAppImage      = "mcr.microsoft.com/dotnet/core/samples:aspnetapp"
	deletePolicy        = meta.DeletePropagationForeground
	targetContainerName = "target"
	namespace           = "default"
	dumperImage         = "kubectl-shovel/dumper-integration-tests"
	tempDirPattern      = "*-kubectl-shovel"
)

type TestCase struct {
	name       string
	args       []string
	pod        *core.Pod
	hostOutput bool
}

func newTestKubeClient() *kubernetes.Client {
	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := k8s.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return &kubernetes.Client{
		Namespace: namespace,
		Clientset: clientset,
	}
}

func setup(t *testing.T, pod *core.Pod) func() {
	t.Helper()
	k := newTestKubeClient()

	fmt.Println("Deploying target pod to cluster...")
	_, err := k.CoreV1().Pods(namespace).Create(
		context.Background(),
		pod,
		meta.CreateOptions{},
	)
	require.NoError(t, err)

	fmt.Println("Waiting app to start...")
	_, err = k.WaitPod(pod.ObjectMeta.Labels)
	require.NoError(t, err)

	return func() {
		_ = k.CoreV1().Pods(namespace).Delete(
			context.TODO(),
			pod.Name,
			meta.DeleteOptions{PropagationPolicy: &deletePolicy},
		)
	}
}

func generateRandomPodMeta() meta.ObjectMeta {
	name := fmt.Sprintf("%s-%s", sampleAppName, uuid.NewString())

	return meta.ObjectMeta{
		Name: name,
		Labels: map[string]string{
			"app": name,
		},
	}
}

func singleContainerPod() *core.Pod {
	return &core.Pod{
		ObjectMeta: generateRandomPodMeta(),
		Spec: core.PodSpec{
			Containers: []core.Container{
				{
					Name:  targetContainerName,
					Image: sampleAppImage,
				},
			},
		},
	}
}

func multiContainerPod() *core.Pod {
	return &core.Pod{
		ObjectMeta: generateRandomPodMeta(),
		Spec: core.PodSpec{
			Containers: []core.Container{
				{
					Name:  targetContainerName,
					Image: sampleAppImage,
				},
				{
					Name:  "sidecar",
					Image: "gcr.io/google_containers/pause-amd64:3.1",
				},
			},
		},
	}
}

func multiContainerPodWithDefaultContainer() *core.Pod {
	objectMeta := generateRandomPodMeta()
	objectMeta.Annotations = map[string]string{
		"kubectl.kubernetes.io/default-container": targetContainerName,
	}
	return &core.Pod{
		ObjectMeta: objectMeta,
		Spec: core.PodSpec{
			Containers: []core.Container{
				{
					Name:  targetContainerName,
					Image: sampleAppImage,
				},
				{
					Name:  "sidecar",
					Image: "gcr.io/google_containers/pause-amd64:3.1",
				},
			},
		},
	}
}

func multiContainerPodWithSharedMount() *core.Pod {
	return &core.Pod{
		ObjectMeta: generateRandomPodMeta(),
		Spec: core.PodSpec{
			Containers: []core.Container{
				{
					Name:  targetContainerName,
					Image: sampleAppImage,
					VolumeMounts: []core.VolumeMount{
						{Name: "shared-path-to-tmp", MountPath: globals.PathTmpFolder},
					},
				},
				{
					Name:  "sidecar",
					Image: "gcr.io/google_containers/pause-amd64:3.1",
					VolumeMounts: []core.VolumeMount{
						{Name: "shared-path-to-tmp", MountPath: globals.PathTmpFolder},
					},
				},
			},
			Volumes: []core.Volume{
				{
					Name: "shared-path-to-tmp",
					VolumeSource: core.VolumeSource{
						EmptyDir: &core.EmptyDirVolumeSource{},
					},
				},
			},
		},
	}
}

func cases(additional ...TestCase) []TestCase {
	basic := []TestCase{
		{
			name: "Basic test",
			args: []string{},
			pod:  singleContainerPod(),
		},
		{
			name:       "Store output on host",
			args:       []string{"store-output-on-host"},
			pod:        singleContainerPod(),
			hostOutput: true,
		},
		{
			name: "MultiContainer pod",
			args: []string{
				"--container",
				targetContainerName,
			},
			pod: multiContainerPod(),
		},
		{
			name: "MultiContainer pod with default-container annotation",
			args: []string{},
			pod:  multiContainerPodWithDefaultContainer(),
		},
		{
			name: "MultiContainer pod with shared mount",
			args: []string{
				"--container",
				targetContainerName,
			},
			pod: multiContainerPodWithSharedMount(),
		},
	}

	return append(basic, additional...)
}
