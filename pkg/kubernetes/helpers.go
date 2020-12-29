package kubernetes

import (
	"fmt"
	"strings"

	v1 "k8s.io/api/core/v1"
)

// ContainerInfo is information about container struct
type ContainerInfo struct {
	Runtime string
	ID      string
}

// JobVolume is helper struct to describe job volume
type JobVolume struct {
	Name      string
	HostPath  string
	MountPath string
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

// NewJobVolume create new helper job volume
func NewJobVolume(containerInfo *ContainerInfo) *JobVolume {
	if containerInfo.Runtime == "containerd" {
		return &JobVolume{
			Name: "tmp",
			HostPath: fmt.Sprintf(
				"/run/containerd/io.containerd.runtime.v2.task/k8s.io/%s/rootfs/tmp",
				containerInfo.ID,
			),
			MountPath: "/tmp",
		}
	}

	return &JobVolume{
		Name:      "dockerfs",
		HostPath:  "/var/lib/docker",
		MountPath: "/var/lib/docker",
	}
}
