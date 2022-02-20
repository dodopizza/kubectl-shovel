package kubernetes

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/wait"
)

type PodInfo struct {
	Name       string
	Node       string
	containers []core.ContainerStatus
}

func NewPodInfo(pod *core.Pod) *PodInfo {
	return &PodInfo{
		Name:       pod.Name,
		Node:       pod.Spec.NodeName,
		containers: pod.Status.ContainerStatuses,
	}
}

// GetContainerNames returns container names associated with pod
func (p *PodInfo) GetContainerNames() []string {
	names := make([]string, len(p.containers))

	for i, cs := range p.containers {
		names[i] = cs.Name
	}

	return names
}

// FindContainerInfo returns container info with specified name or error
func (p *PodInfo) FindContainerInfo(container string) (*ContainerInfo, error) {
	if container == "" && len(p.containers) > 1 {
		return nil, fmt.Errorf(
			"container name must be specified for pod %s, choose one of: [%s]",
			p.Name,
			strings.Join(p.GetContainerNames(), " "),
		)
	}

	cs, err := p.FindContainerStatus(container)
	if err != nil {
		return nil, err
	}
	return NewContainerInfo(cs), nil
}

// FindContainerStatus returns container status info for specified container
func (p *PodInfo) FindContainerStatus(container string) (*core.ContainerStatus, error) {
	if container == "" {
		return &p.containers[0], nil
	}

	for _, cs := range p.containers {
		if cs.Name == container {
			return &cs, nil
		}
	}

	return nil, fmt.Errorf("container %s is not valid for pod %s", container, p.Name)
}

// GetPodInfo get info about pod by name
func (k8s *Client) GetPodInfo(name string) (*PodInfo, error) {
	pod, err := k8s.
		CoreV1().
		Pods(k8s.Namespace).
		Get(context.Background(), name, meta.GetOptions{})

	if err != nil {
		return nil, err
	}

	return NewPodInfo(pod), nil
}

// WaitPod will wait until pod to start or failed
func (k8s *Client) WaitPod(labelSelector map[string]string) (*PodInfo, error) {
	var pod *core.Pod

	err := wait.Poll(1*time.Second, 5*time.Minute,
		func() (bool, error) {
			podList, err := k8s.
				CoreV1().
				Pods(k8s.Namespace).
				List(
					context.Background(),
					meta.ListOptions{
						LabelSelector: labels.Set(labelSelector).String(),
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
			case core.PodFailed:
				return false, errors.New("pod has been failed")
			case core.PodSucceeded, core.PodRunning:
				return true, nil
			default:
				return false, nil
			}
		},
	)
	if err != nil {
		return nil, err
	}

	return NewPodInfo(pod), nil
}

// ReadPodLogs stream logs from pod
func (k8s *Client) ReadPodLogs(podName, containerName string) (io.ReadCloser, error) {
	req := k8s.CoreV1().
		Pods(k8s.Namespace).
		GetLogs(podName, &core.PodLogOptions{
			Container: containerName,
			Follow:    true,
		})

	return req.Stream(context.Background())
}
