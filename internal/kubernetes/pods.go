package kubernetes

import (
	"context"
	"fmt"
	"io"
	"time"

	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
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

// WaitPod will wait pod to start
func (k8s *Client) WaitPod(jobName string) (string, error) {
	var pod *v1.Pod
	err := wait.Poll(1*time.Second, 5*time.Minute,
		func() (bool, error) {
			podList, err := k8s.
				CoreV1().
				Pods(k8s.Namespace).
				List(
					context.Background(),
					metav1.ListOptions{
						LabelSelector: fmt.Sprintf("job-name=%s", jobName),
					},
				)
			if err != nil {
				return false, err
			}

			if len(podList.Items) == 0 {
				return false, nil
			}

			pod = &podList.Items[0]
			switch pod.Status.Phase {
			case v1.PodFailed:
				return false, fmt.Errorf("pod failed")
			case v1.PodSucceeded, v1.PodRunning:
				return true, nil
			default:
				return false, nil
			}
		},
	)
	if err != nil {
		return "", err
	}

	return pod.Name, nil
}

// ReadPodLogs stream logs from pod
func (k8s *Client) ReadPodLogs(name string) (io.ReadCloser, error) {
	req := k8s.CoreV1().
		Pods(k8s.Namespace).
		GetLogs(name, &apiv1.PodLogOptions{
			Follow: true,
		})

	readCloser, err := req.Stream(context.Background())
	if err != nil {
		return nil, err
	}

	return readCloser, nil
}
