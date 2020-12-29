package kubernetes

import (
	"strings"

	v1 "k8s.io/api/core/v1"
)

// ContainerInfo is information about container struct
type ContainerInfo struct {
	Runtime string
	ID      string
}

func int32Ptr(i int32) *int32 {
	return &i
}

// GetContainerInfo helps to get info about container
func GetContainerInfo(pod *v1.Pod) *ContainerInfo {
	containerInfo := strings.Split(pod.Status.ContainerStatuses[0].ContainerID, "://")
	return &ContainerInfo{
		Runtime: containerInfo[0],
		ID:      containerInfo[1],
	}
}
