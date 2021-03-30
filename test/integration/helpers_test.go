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
	sampleAppName  = "sample-app"
	sampleAppImage = "mcr.microsoft.com/dotnet/core/samples:aspnetapp"
	deletePolicy   = metav1.DeletePropagationForeground
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

func setup(t *testing.T) (string, func()) {
	t.Helper()
	targetPodName := fmt.Sprintf(
		"%s-%s",
		sampleAppName,
		newRandom(),
	)
	k8s := newTestKubeClient()
	labels := map[string]string{
		"app": targetPodName,
	}

	sampleAppPod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:   targetPodName,
			Labels: labels,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "target",
					Image: sampleAppImage,
				},
				{
					Name:  "sidecar",
					Image: "gcr.io/google_containers/pause-amd64:3.1",
				},
			},
		},
	}

	fmt.Println("Deploying dotnet sample app to cluster...")
	_, err := k8s.CoreV1().Pods(namespace).Create(
		context.Background(),
		sampleAppPod,
		metav1.CreateOptions{},
	)
	require.NoError(t, err)

	fmt.Println("Waiting app to start...")
	_, err = k8s.WaitPod(labels)
	require.NoError(t, err)

	return targetPodName, func() {
		_ = k8s.CoreV1().Pods(namespace).Delete(
			context.TODO(),
			targetPodName,
			metav1.DeleteOptions{PropagationPolicy: &deletePolicy},
		)
	}
}

func newRandom() string {
	return uuid.New().String()
}
