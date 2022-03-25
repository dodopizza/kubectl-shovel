package kubernetes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	core "k8s.io/api/core/v1"

	"github.com/dodopizza/kubectl-shovel/internal/globals"
)

// ContainerInfo is an information about container struct
type ContainerInfo struct {
	ID      string
	Name    string
	Runtime string
}

// ContainerConfigInfo is an information about container state located at node
type ContainerConfigInfo struct {
	ID            string
	HostProcessID int
	Mounts        []ContainerMount
	RootFS        string
	Runtime       string
}

// ContainerMount is an information about available container mounts
type ContainerMount struct {
	Source      string
	Destination string
}

// NewContainerInfo returns container info mapped from core.ContainerStatus
func NewContainerInfo(cs *core.ContainerStatus) *ContainerInfo {
	containerInfo := strings.Split(cs.ContainerID, "://")

	return &ContainerInfo{
		ID:      containerInfo[1],
		Name:    cs.Name,
		Runtime: containerInfo[0],
	}
}

// NewContainerConfigInfo returns container configuration state
func NewContainerConfigInfo(runtime, id string) (*ContainerConfigInfo, error) {
	switch runtime {
	case "containerd":
		return containerd(id)
	case "docker":
		return docker(id)
	default:
		return nil, fmt.Errorf("unsupported runtime: %s for container: %s", runtime, id)
	}
}

// GetTmpSource returns mount point for /tmp folder, depending upon container runtime and existing mounts
// If container contains mounts to /tmp folder, this mount source path will be used, otherwise – container rootfs
func (c *ContainerConfigInfo) GetTmpSource() string {
	for _, mount := range c.Mounts {
		if mount.Destination == globals.PathTmpFolder {
			return mount.Source
		}
	}

	return fmt.Sprintf("%s%s", c.RootFS, globals.PathTmpFolder)
}

// GetContainerFSVolume returns JobVolume (mounted from host) that contains container definitions,
// depending upon container runtime
func (c *ContainerInfo) GetContainerFSVolume() JobVolume {
	if c.Runtime == "containerd" {
		return JobVolume{
			Name:      "containerdfs",
			HostPath:  globals.PathContainerDFS,
			MountPath: globals.PathContainerDFS,
		}
	}

	return JobVolume{
		Name:      "dockerfs",
		HostPath:  globals.PathDockerFS,
		MountPath: globals.PathDockerFS,
	}
}

// GetContainerSharedVolumes returns JobVolume (mounted from host) that contains container additional mounts,
// depending upon container runtime
func (c *ContainerInfo) GetContainerSharedVolumes() JobVolume {
	if c.Runtime == "containerd" {
		return JobVolume{
			Name:      "containerdvolumes",
			HostPath:  globals.PathContainerDVolumes,
			MountPath: globals.PathContainerDVolumes,
		}
	}

	return JobVolume{
		Name:      "dockervolumes",
		HostPath:  globals.PathDockerVolumes,
		MountPath: globals.PathDockerVolumes,
	}
}

func docker(id string) (*ContainerConfigInfo, error) {
	mountFile := fmt.Sprintf("%s/image/overlay2/layerdb/mounts/%s/mount-id", globals.PathDockerFS, id)
	mountId, err := ioutil.ReadFile(mountFile)
	if err != nil {
		return nil, err
	}

	stateFile, err := os.Open(fmt.Sprintf("%s/containers/%s/config.v2.json", globals.PathDockerFS, id))
	if err != nil {
		return nil, err
	}

	state := &struct {
		State struct {
			Pid int `json:"Pid"`
		} `json:"State"`
		MountPoints map[string]struct {
			Source      string `json:"Source"`
			Destination string `json:"Destination"`
		} `json:"MountPoints"`
	}{}
	if err := json.NewDecoder(stateFile).Decode(state); err != nil {
		return nil, err
	}

	mounts := make([]ContainerMount, 0)
	for _, mount := range state.MountPoints {
		mounts = append(mounts, ContainerMount{
			Source:      mount.Source,
			Destination: mount.Destination,
		})
	}

	return &ContainerConfigInfo{
		ID:            id,
		HostProcessID: state.State.Pid,
		Mounts:        mounts,
		Runtime:       "docker",
		RootFS:        fmt.Sprintf("%s/overlay2/%s/merged", globals.PathDockerFS, mountId),
	}, nil
}

func containerd(id string) (*ContainerConfigInfo, error) {
	file, err := os.Open(fmt.Sprintf("%s/runc/k8s.io/%s/state.json", globals.PathContainerDFS, id))
	if err != nil {
		return nil, err
	}

	state := &struct {
		InitProcessPID int `json:"init_process_pid"`
		Config         struct {
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

	mounts := make([]ContainerMount, len(state.Config.Mounts))
	for i, mount := range state.Config.Mounts {
		mounts[i] = ContainerMount{
			Source:      mount.Source,
			Destination: mount.Destination,
		}
	}

	return &ContainerConfigInfo{
		ID:            id,
		HostProcessID: state.InitProcessPID,
		Mounts:        mounts,
		Runtime:       "containerd",
		RootFS:        state.Config.RootFS,
	}, nil
}
