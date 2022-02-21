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
	ID      string
	Runtime string
}

type containerConfig struct {
	RootFS string
	Mounts []containerMount
}

type containerMount struct {
	Source      string
	Destination string
}

// NewContainerInfo returns container info mapped from core.ContainerStatus
func NewContainerInfo(cs *core.ContainerStatus) *ContainerInfo {
	containerInfo := strings.Split(cs.ContainerID, "://")

	return &ContainerInfo{
		ID:      containerInfo[1],
		Runtime: containerInfo[0],
	}
}

// NewContainerInfoRaw returns container info with specified runtime and id
func NewContainerInfoRaw(runtime, id string) *ContainerInfo {
	return &ContainerInfo{
		Runtime: runtime,
		ID:      id,
	}
}

// GetTmpSource returns mount point for /tmp folder, depending upon container runtime and existing mounts
// If container contains mounts to /tmp folder, this mount source path will be used, otherwise – container rootfs
func (c *ContainerInfo) GetTmpSource() (string, error) {
	config, err := c.config()
	if err != nil {
		return "", nil
	}

	for _, mount := range config.Mounts {
		if mount.Destination == "/tmp" {
			return mount.Source, nil
		}
	}

	return fmt.Sprintf("%s/tmp", config.RootFS), nil
}

// GetContainerFSVolume returns JobVolume (mounted from host) that contains container definitions,
// depending upon container runtime
func (c *ContainerInfo) GetContainerFSVolume() JobVolume {
	if c.containerd() {
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
	if c.containerd() {
		return JobVolume{
			Name:      "containerdvolumes",
			HostPath:  "/var/lib/kubelet/pods",
			MountPath: "/var/lib/kubelet/pods",
		}
	}

	return JobVolume{
		Name:      "dockervolumes",
		HostPath:  "/var/lib/docker",
		MountPath: "/var/lib/docker",
	}
}

func (c *ContainerInfo) config() (*containerConfig, error) {
	if c.containerd() {
		return c.containerdConfig()
	}
	return c.dockerConfig()
}

func (c *ContainerInfo) dockerConfig() (*containerConfig, error) {
	path := "/var/lib/docker"

	mountFile := fmt.Sprintf("%s/image/overlay2/layerdb/mounts/%s/mount-id", path, c.ID)
	mountId, err := ioutil.ReadFile(mountFile)
	if err != nil {
		return nil, err
	}

	stateFile, err := os.Open(fmt.Sprintf("%s/containers/%s/config.v2.json", path, c.ID))
	if err != nil {
		return nil, err
	}

	state := &struct {
		MountPoints map[string]struct {
			Source      string `json:"Source"`
			Destination string `json:"Destination"`
		} `json:"MountPoints"`
	}{}
	if err := json.NewDecoder(stateFile).Decode(state); err != nil {
		return nil, err
	}

	mounts := make([]containerMount, 0)
	for _, mount := range state.MountPoints {
		mounts = append(mounts, containerMount{
			Source:      mount.Source,
			Destination: mount.Destination,
		})
	}

	return &containerConfig{
		RootFS: fmt.Sprintf("%s/overlay2/%s/merged", path, mountId),
		Mounts: mounts,
	}, nil
}

func (c *ContainerInfo) containerdConfig() (*containerConfig, error) {
	file, err := os.Open(fmt.Sprintf("/run/containerd/runc/k8s.io/%s/state.json", c.ID))
	if err != nil {
		return nil, err
	}

	state := &struct {
		Config struct {
			RootFS string `json:"rootfs"`
			Mounts []struct {
				Source      string `json:"source"`
				Destination string `json:"destination"`
			} `json:"mounts"`
		} `json:"config"`
	}{}

	if err := json.NewDecoder(file).Decode(state); err != nil {
		return nil, err
	}

	mounts := make([]containerMount, len(state.Config.Mounts))
	for i, mount := range state.Config.Mounts {
		mounts[i] = containerMount{
			Source:      mount.Source,
			Destination: mount.Destination,
		}
	}

	return &containerConfig{
		RootFS: state.Config.RootFS,
		Mounts: mounts,
	}, nil
}

func (c *ContainerInfo) containerd() bool {
	return c.Runtime == "containerd"
}
