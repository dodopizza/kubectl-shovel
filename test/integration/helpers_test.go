// +build integration

package integration_test

import (
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	"github.com/dodopizza/kubectl-shovel/internal/kubernetes"
)

var deletePolicy = metav1.DeletePropagationForeground

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
