package kubernetes

import (
	"encoding/json"
	"errors"
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
	switch ci.Runtime {
	case "docker":
		return ci.dockerMountPoint()
	case "containerd":
		return ci.containerDMountPoint()
	default:
		return "", errors.New("unknown container runtime")
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
