package kubernetes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	core "k8s.io/api/core/v1"
	"os"
	"strings"
)

// ContainerInfo is information about container struct
type ContainerInfo struct {
	Runtime string
	ID      string
}

// NewContainerInfo returns container info mapped from core.ContainerStatus
func NewContainerInfo(cs *core.ContainerStatus) *ContainerInfo {
	containerInfo := strings.Split(cs.ContainerID, "://")

	return &ContainerInfo{
		Runtime: containerInfo[0],
		ID:      containerInfo[1],
	}
}

// GetMountPoint returns mount point depending ContainerRuntime
func (c *ContainerInfo) GetMountPoint() (string, error) {
	if c.Runtime == "docker" {
		return c.dockerMountPoint()
	}
	return c.containerDMountPoint()
}

// GetContainerFSVolume returns JobVolume (mounted from host) that contains container definitions,
// depending upon container runtime
func (c *ContainerInfo) GetContainerFSVolume() JobVolume {
	if c.Runtime == "containerd" {
		return JobVolume{
			Name:      "containerdfs",
			HostPath:  "/run/containerd",
			MountPath: "/run/containerd",
		}
	}

	return JobVolume{
		Name:      "dockerfs",
		HostPath:  "/var/lib/docker",
		MountPath: "/var/lib/docker",
	}
}

// GetContainerSharedVolumes returns JobVolume (mounted from host) that contains container additional mounts,
// depending upon container runtime
func (c *ContainerInfo) GetContainerSharedVolumes() JobVolume {
	if c.Runtime == "containerd" {
		return JobVolume{
			Name:      "containerdvolumes",
			HostPath:  "/var/lib/kubelet/pods",
			MountPath: "/var/lib/kubelet/pods",
		}
	}

	return JobVolume{
		Name:      "dockervolumes",
		HostPath:  "/",
		MountPath: "/",
	}
}

func (c *ContainerInfo) dockerMountPoint() (string, error) {
	id, err := ioutil.ReadFile(
		fmt.Sprintf(
			"/var/lib/docker/image/overlay2/layerdb/mounts/%s/mount-id",
			c.ID,
		),
	)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"/var/lib/docker/overlay2/%s/merged",
		string(id),
	), nil
}

func (c *ContainerInfo) containerDMountPoint() (string, error) {
	file, err := os.Open(
		fmt.Sprintf(
			"/run/containerd/runc/k8s.io/%s/state.json",
			c.ID,
		),
	)
	if err != nil {
		return "", err
	}
	state := &struct {
		Config struct {
			RootFS string `json:"rootfs"`
		} `json:"config"`
	}{}
	if err := json.NewDecoder(file).Decode(state); err != nil {
		return "", err
	}

	return state.Config.RootFS, nil
}
