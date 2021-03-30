// +build integration

package integration_test

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/dodopizza/kubectl-shovel/internal/kubernetes"
)

var (
	sampleAppName       = "sample-app"
	sampleAppImage      = "mcr.microsoft.com/dotnet/core/samples:aspnetapp"
	deletePolicy        = metav1.DeletePropagationForeground
	targetContainerName = "target"
	namespace           = "default"
	dumperImage         = "kubectl-shovel/dumper-integration-tests"
	tempDirPattern      = "*-kubectl-shovel"
)

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

func setup(t *testing.T, pod *v1.Pod) func() {
	t.Helper()
	k8s := newTestKubeClient()

	fmt.Println("Deploying target pod to cluster...")
	_, err := k8s.CoreV1().Pods(namespace).Create(
		context.Background(),
		pod,
		metav1.CreateOptions{},
	)
	require.NoError(t, err)

	fmt.Println("Waiting app to start...")
	_, err = k8s.WaitPod(map[string]string{
		"app": pod.Name,
	})
	require.NoError(t, err)

	return func() {
		_ = k8s.CoreV1().Pods(namespace).Delete(
			context.TODO(),
			pod.Name,
			metav1.DeleteOptions{PropagationPolicy: &deletePolicy},
		)
	}
}

func newRandomString() string {
	return uuid.NewString()
}

func sampleAppPod() *v1.Pod {
	targetPodName := fmt.Sprintf(
		"%s-%s",
		sampleAppName,
		newRandomString(),
	)
	labels := map[string]string{
		"app": targetPodName,
	}
	return &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:   targetPodName,
			Labels: labels,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  targetContainerName,
					Image: sampleAppImage,
				},
			},
		},
	}
}

func multiContainerPod() *v1.Pod {
	targetPodName := fmt.Sprintf(
		"%s-%s",
		sampleAppName,
		newRandomString(),
	)
	labels := map[string]string{
		"app": targetPodName,
	}
	return &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:   targetPodName,
			Labels: labels,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
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
