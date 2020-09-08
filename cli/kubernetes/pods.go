package kubernetes

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetPodInfo get info about pod by name
func (k8s *Client) GetPodInfo(podName string) (*v1.Pod, error) {
	return k8s.
		CoreV1().
		Pods(k8s.Namespace).
		Get(
			context.Background(),
			podName,
			metav1.GetOptions{},
		)
}
