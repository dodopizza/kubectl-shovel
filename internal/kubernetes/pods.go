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
	Annotations       map[string]string
	Name              string
	Node              string
	containers        []core.Container
	containerStatuses []core.ContainerStatus
}

// NewPodInfo returns PodInfo generated from core.Pod spec
func NewPodInfo(pod *core.Pod) *PodInfo {
	return &PodInfo{
		Annotations:       pod.Annotations,
		Name:              pod.Name,
		Node:              pod.Spec.NodeName,
		containers:        pod.Spec.Containers,
		containerStatuses: pod.Status.ContainerStatuses,
	}
}

// GetContainerNames returns container names associated with pod
func (p *PodInfo) GetContainerNames() []string {
	names := make([]string, len(p.containerStatuses))

	for i, cs := range p.containerStatuses {
		names[i] = cs.Name
	}

	return names
}

// FindContainerInfo returns container info with specified name or error
func (p *PodInfo) FindContainerInfo(container string) (*ContainerInfo, error) {
	_, cs, err := p.findContainerInfo(container)
	if err != nil {
		return nil, err
	}
	return NewContainerInfo(cs), nil
}

// ContainsMountedTmp returns true if container has /tmp folder mounted from host or shared with other container
func (p *PodInfo) ContainsMountedTmp(container string) bool {
	c, _, err := p.findContainerInfo(container)
	if err != nil {
		return false
	}

	for _, mount := range c.VolumeMounts {
		if mount.MountPath == "/tmp" {
			return true
		}
	}

	return false
}

func (p *PodInfo) findContainerInfo(name string) (container *core.Container, status *core.ContainerStatus, err error) {
	count := len(p.containerStatuses)

	// check against multiple containers
	if count > 1 && name == "" {
		err = fmt.Errorf(
			"container name must be specified for pod %s, choose one of: [%s]",
			p.Name,
			strings.Join(p.GetContainerNames(), " "),
		)
		return
	}

	if count == 1 && name == "" {
		container = &p.containers[0]
		status = &p.containerStatuses[0]
		return
	}

	for i := range p.containers {
		if p.containers[i].Name == name {
			container = &p.containers[i]
		}
	}

	for i := range p.containerStatuses {
		if p.containerStatuses[i].Name == name {
			status = &p.containerStatuses[i]
		}
	}

	// not found
	if container == nil || status == nil {
		err = fmt.Errorf("container %s is not valid for pod %s", container, p.Name)
	}

	return
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
