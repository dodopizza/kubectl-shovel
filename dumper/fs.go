package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func mapContainerTmp(containerID string) error {
	containerFS, err := getContainerMountpoint(containerID)
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

func getContainerMountpoint(containerID string) (string, error) {
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
