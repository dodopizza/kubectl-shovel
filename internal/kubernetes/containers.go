package kubernetes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// ContainerInfo is information about container struct
type ContainerInfo struct {
	Runtime string
	ID      string
}

// GetMountPoint returns mount point depending ContainerRuntime
func (ci *ContainerInfo) GetMountPoint() (string, error) {
	if ci.Runtime == "docker" {
		return ci.dockerMountPoint()
	}
	return ci.containerDMountPoint()
}

// GetJobVolume returns helper job volume
func (ci *ContainerInfo) GetJobVolume() *JobVolume {
	if ci.Runtime == "containerd" {
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

func (ci *ContainerInfo) dockerMountPoint() (string, error) {
	id, err := ioutil.ReadFile(
		fmt.Sprintf(
			"/var/lib/docker/image/overlay2/layerdb/mounts/%s/mount-id",
			ci.ID,
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

func (ci *ContainerInfo) containerDMountPoint() (string, error) {
	file, err := os.Open(
		fmt.Sprintf(
			"/run/containerd/runc/k8s.io/%s/state.json",
			ci.ID,
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
