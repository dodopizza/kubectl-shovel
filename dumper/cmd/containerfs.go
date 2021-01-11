package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/dodopizza/kubectl-shovel/internal/kubernetes"
)

func mapContainerTmp(containerInfo kubernetes.ContainerInfo) error {
	var containerFS string
	var err error
	switch containerInfo.Runtime {
	case "docker":
		containerFS, err = getDockerContainerMountpoint(containerInfo.ID)
	case "containerd":
		containerFS, err = getContainerDContainerMountpoint(containerInfo.ID)
	default:
		return errors.New("Unknown container runtime")
	}
	if err != nil {
		return err
	}
	err = os.RemoveAll("/tmp")
	if err != nil {
		return err
	}

	return os.Symlink(
		filepath.Join(containerFS, "tmp"),
		"/tmp",
	)
}

func getDockerContainerMountpoint(containerID string) (string, error) {
	id, err := ioutil.ReadFile(
		fmt.Sprintf(
			"/var/lib/docker/image/overlay2/layerdb/mounts/%s/mount-id",
			containerID,
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

func getContainerDContainerMountpoint(containerID string) (string, error) {
	file, err := os.Open(
		fmt.Sprintf(
			"/run/containerd/runc/k8s.io/%s/state.json",
			containerID,
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
