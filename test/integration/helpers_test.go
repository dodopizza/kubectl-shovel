//go:build integration
// +build integration

package integration_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	k8s "k8s.io/client-go/kubernetes"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/dodopizza/kubectl-shovel/internal/flags"
	"github.com/dodopizza/kubectl-shovel/internal/globals"
	"github.com/dodopizza/kubectl-shovel/internal/kubernetes"
)

var (
	dumperImage = "kubectl-shovel/dumper-integration-tests"

	namespace             = "default"
	sidecarContainerImage = "gcr.io/google_containers/pause:3.1"
	sidecarContainerName  = "sidecar"
	targetPodNamePrefix   = "sample-app"
	targetContainerImage  = "kubectl-shovel/sample-integration-tests"
	targetContainerName   = "target"
)

type TestCase struct {
	name       string
	args       map[string]string
	pod        *core.Pod
	output     string
	hostOutput bool
}

func (tc *TestCase) FormatArgs(command string) []string {
	args := flags.NewArgs().
		AppendRaw(command).
		Append("pod-name", tc.pod.Name).
		Append("image", dumperImage)

	if tc.hostOutput {
		args.AppendKey("store-output-on-host")
	} else {
		args.Append("output", tc.output)
	}

	for key, value := range tc.args {
		args.Append(key, value)
	}

	return args.Get()
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

func setup(t *testing.T, tc *TestCase, prefix string) func() {
	t.Parallel()
	t.Helper()
	k := newTestKubeClient()

	t.Log("Deploying target pod to cluster...")
	_, err := k.CoreV1().Pods(namespace).Create(
		context.Background(),
		tc.pod,
		meta.CreateOptions{},
	)
	require.NoError(t, err)

	t.Log("Waiting target pod to start...")
	_, err = k.WaitPod(tc.pod.ObjectMeta.Labels)
	require.NoError(t, err)

	if !tc.hostOutput {
		dir, _ := ioutil.TempDir("", "*-kubectl-shovel")
		tc.output = filepath.Join(dir, prefix)
		t.Logf("Output for test case will be stored at: %s\n", tc.output)
	}

	return func() {
		t.Logf("Delete test pod: %s\n", tc.pod.Name)

		policy := meta.DeletePropagationForeground
		_ = k.CoreV1().Pods(namespace).Delete(
			context.TODO(),
			tc.pod.Name,
			meta.DeleteOptions{PropagationPolicy: &policy},
		)

		if !tc.hostOutput {
			dir := filepath.Dir(tc.output)
			t.Logf("Cleanup test case output dir: %s\n", dir)
			_ = os.Remove(dir)
		}
	}
}

func generateRandomPodMeta() meta.ObjectMeta {
	name := fmt.Sprintf("%s-%s", targetPodNamePrefix, uuid.NewString())

	return meta.ObjectMeta{
		Name: name,
		Labels: map[string]string{
			"app": name,
		},
	}
}

func targetContainer() core.Container {
	return core.Container{
		Name:  targetContainerName,
		Image: targetContainerImage,
		Ports: []core.ContainerPort{{
			ContainerPort: 6000,
			Name:          "app",
			Protocol:      "TCP",
		}},
		LivenessProbe: &core.Probe{
			ProbeHandler: core.ProbeHandler{
				HTTPGet: &core.HTTPGetAction{
					Path: "/health/live",
					Port: intstr.IntOrString{
						Type:   intstr.String,
						StrVal: "app",
					},
					Scheme: "http",
				},
			},
			InitialDelaySeconds: 2,
			TimeoutSeconds:      1,
			PeriodSeconds:       1,
			SuccessThreshold:    1,
			FailureThreshold:    5,
		},
		TerminationMessagePolicy: core.TerminationMessageFallbackToLogsOnError,
	}
}

func sidecarContainer() core.Container {
	return core.Container{
		Name:  sidecarContainerName,
		Image: sidecarContainerImage,
	}
}

func singleContainerPod() *core.Pod {
	return &core.Pod{
		ObjectMeta: generateRandomPodMeta(),
		Spec: core.PodSpec{
			Containers: []core.Container{targetContainer()},
		},
	}
}

func multiContainerPod() *core.Pod {
	return &core.Pod{
		ObjectMeta: generateRandomPodMeta(),
		Spec: core.PodSpec{
			Containers: []core.Container{targetContainer(), sidecarContainer()},
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
			Containers: []core.Container{targetContainer(), sidecarContainer()},
		},
	}
}

func multiContainerPodWithSharedMount() *core.Pod {
	volumes := []core.Volume{
		{
			Name: "shared-path-to-tmp",
			VolumeSource: core.VolumeSource{
				EmptyDir: &core.EmptyDirVolumeSource{},
			},
		},
	}
	mounts := []core.VolumeMount{
		{
			Name:      "shared-path-to-tmp",
			MountPath: globals.PathTmpFolder,
		},
	}

	sidecar := sidecarContainer()
	sidecar.VolumeMounts = mounts

	target := targetContainer()
	target.VolumeMounts = mounts

	return &core.Pod{
		ObjectMeta: generateRandomPodMeta(),
		Spec: core.PodSpec{
			Containers: []core.Container{target, sidecar},
			Volumes:    volumes,
		},
	}
}

func cases(additional ...TestCase) []TestCase {
	basic := []TestCase{
		{
			name: "Basic test",
			args: map[string]string{},
			pod:  singleContainerPod(),
		},
		{
			name:       "Store output on host",
			args:       map[string]string{},
			pod:        singleContainerPod(),
			hostOutput: true,
		},
		{
			name: "MultiContainer pod",
			args: map[string]string{
				"container": targetContainerName,
			},
			pod: multiContainerPod(),
		},
		{
			name: "MultiContainer pod with default-container annotation",
			args: map[string]string{},
			pod:  multiContainerPodWithDefaultContainer(),
		},
		{
			name: "MultiContainer pod with shared mount",
			args: map[string]string{
				"container": targetContainerName,
			},
			pod: multiContainerPodWithSharedMount(),
		},
	}

	return append(basic, additional...)
}
