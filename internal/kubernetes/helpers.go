package kubernetes

import (
	"fmt"
	"strings"

	v1 "k8s.io/api/core/v1"
)

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
func GetContainerInfo(
	pod *v1.Pod,
	containerName string,
) (*ContainerInfo, error) {
	if containerName == "" && len(pod.Status.ContainerStatuses) > 1 {
		return nil, fmt.Errorf(
			"Container name must be specified for pod %s, choose one of: [%s]",
			pod.Name,
			strings.Join(getContainerNames(pod), " "),
		)
	}
	var cs v1.ContainerStatus
	if containerName != "" {
		var err error
		cs, err = getContainerInfoByName(pod, containerName)
		if err != nil {
			return nil, err
		}
	} else {
		cs = pod.Status.ContainerStatuses[0]
	}

	containerInfo := strings.Split(cs.ContainerID, "://")
	return &ContainerInfo{
		Runtime: containerInfo[0],
		ID:      containerInfo[1],
	}, nil
}

func getContainerNames(pod *v1.Pod) []string {
	names := make([]string, len(pod.Status.ContainerStatuses))
	for i, cs := range pod.Status.ContainerStatuses {
		names[i] = cs.Name
	}

	return names
}

func getContainerInfoByName(
	pod *v1.Pod,
	containerName string,
) (v1.ContainerStatus, error) {
	for _, cs := range pod.Status.ContainerStatuses {
		if cs.Name == containerName {
			return cs, nil
		}
	}

	return v1.ContainerStatus{}, fmt.Errorf(
		"Container %s is not valid for pod %s",
		containerName,
		pod.Name,
	)
}

// NewJobVolume create new helper job volume
func NewJobVolume(containerInfo *ContainerInfo) *JobVolume {
	if containerInfo.Runtime == "containerd" {
		return &JobVolume{
			Name:      "containerdfs",
			HostPath:  "/run/containerd",
			MountPath: "/run/containerd",
		}
	}

	return &JobVolume{
		Name:      "dockerfs",
		HostPath:  "/var/lib/docker",
		MountPath: "/var/lib/docker",
	}
}
