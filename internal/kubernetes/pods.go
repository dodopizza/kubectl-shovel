package kubernetes

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/dodopizza/kubectl-shovel/internal/globals"

	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/wait"
)

type PodInfo struct {
	Annotations             map[string]string
	Name                    string
	Namespace               string
	Node                    string
	containers              []core.Container
	containerStatuses       []core.ContainerStatus
	initContainers          []core.Container
	initContainerStatuses   []core.ContainerStatus
}

// NewPodInfo returns PodInfo generated from core.Pod spec
func NewPodInfo(pod *core.Pod) *PodInfo {
	return &PodInfo{
		Annotations:           pod.Annotations,
		Name:                  pod.Name,
		Namespace:             pod.Namespace,
		Node:                  pod.Spec.NodeName,
		containers:            pod.Spec.Containers,
		containerStatuses:     pod.Status.ContainerStatuses,
		initContainers:        pod.Spec.InitContainers,
		initContainerStatuses: pod.Status.InitContainerStatuses,
	}
}

// GetContainerNames returns container names associated with pod
func (p *PodInfo) GetContainerNames() []string {
	names := make([]string, len(p.containerStatuses)+len(p.initContainerStatuses))

	for i, cs := range p.containerStatuses {
		names[i] = cs.Name
	}

	offset := len(p.containerStatuses)
	for i, cs := range p.initContainerStatuses {
		names[offset+i] = cs.Name
	}

	return names
}

// FindContainerInfo returns container info with specified name or error
func (p *PodInfo) FindContainerInfo(container string) (*ContainerInfo, error) {
	_, cs, err := p.findContainerInfo(container)
	if err != nil {
		return nil, err
	}
	
	// Check if container has a valid container ID
	if cs.ContainerID == "" {
		return nil, fmt.Errorf("container '%s' does not have a valid container ID. Container may not be running or has already completed", container)
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
		if mount.MountPath == globals.PathTmpFolder {
			return true
		}
	}

	return false
}

// IsInitContainer returns true if the specified container is an init container
func (p *PodInfo) IsInitContainer(containerName string) bool {
	for _, c := range p.initContainers {
		if c.Name == containerName {
			return true
		}
	}
	return false
}

func (p *PodInfo) findContainerInfo(name string) (container *core.Container, status *core.ContainerStatus, err error) {
	regularCount := len(p.containerStatuses)
	totalCount := regularCount + len(p.initContainerStatuses)

	// check against multiple containers
	if totalCount > 1 && name == "" {
		err = fmt.Errorf(
			"container name must be specified for pod %s, choose one of: [%s]",
			p.Name,
			strings.Join(p.GetContainerNames(), " "),
		)
		return
	}

	if regularCount == 1 && name == "" {
		container = &p.containers[0]
		status = &p.containerStatuses[0]
		return
	}

	// First search in regular containers
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

	// If container or status is not found in regular containers, search in init containers
	if container == nil || status == nil {
		for i := range p.initContainers {
			if p.initContainers[i].Name == name {
				container = &p.initContainers[i]
			}
		}

		for i := range p.initContainerStatuses {
			if p.initContainerStatuses[i].Name == name {
				status = &p.initContainerStatuses[i]
			}
		}
	}

	// not found in either regular or init containers
	if container == nil || status == nil {
		err = fmt.Errorf("container '%s' is not valid for pod %s", name, p.Name)
	}

	return
}

// GetPodInfo get info about pod by name
func (k *Client) GetPodInfo(name string) (*PodInfo, error) {
	pod, err := k.
		CoreV1().
		Pods(k.Namespace).
		Get(context.Background(), name, meta.GetOptions{})

	if err != nil {
		return nil, err
	}

	return NewPodInfo(pod), nil
}

// WaitPod will wait until pod to start or failed
func (k *Client) WaitPod(labelSelector map[string]string) (*PodInfo, error) {
	var pod *core.Pod

	err := wait.Poll(1*time.Second, 5*time.Minute, func() (bool, error) {
		options := meta.ListOptions{
			LabelSelector: labels.Set(labelSelector).String(),
		}
		list, err := k.CoreV1().
			Pods(k.Namespace).
			List(context.Background(), options)

		if err != nil {
			return false, err
		}

		if len(list.Items) == 0 {
			return false, nil
		}

		pod = &list.Items[0]
		switch pod.Status.Phase {
		case core.PodFailed:
			message := ""

			for _, c := range pod.Status.ContainerStatuses {
				if c.State.Terminated != nil {
					state := c.State.Terminated
					message += fmt.Sprintf("container: %s, reason: %s, state:\n%s",
						c.Name,
						state.Reason,
						state.Message)
				}
				if c.State.Waiting != nil {
					state := c.State.Waiting
					message += fmt.Sprintf("container: %s, reason: %s, state:\n%s",
						c.Name,
						state.Reason,
						state.Message)
				}
			}

			return false, fmt.Errorf("pod has been failed, container statuses:\n%s", message)
		case core.PodSucceeded, core.PodRunning:
			return true, nil
		default:
			return false, nil
		}
	})

	if err != nil {
		return nil, err
	}

	return NewPodInfo(pod), nil
}

// ReadPodLogs stream logs from pod
func (k *Client) ReadPodLogs(podName, containerName string) (io.ReadCloser, error) {
	req := k.CoreV1().
		Pods(k.Namespace).
		GetLogs(podName, &core.PodLogOptions{
			Container: containerName,
			Follow:    true,
		})

	return req.Stream(context.Background())
}
