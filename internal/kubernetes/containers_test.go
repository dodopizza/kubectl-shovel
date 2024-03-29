package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/require"
	core "k8s.io/api/core/v1"

	"github.com/dodopizza/kubectl-shovel/internal/globals"
)

func Test_NewContainerInfo(t *testing.T) {
	testCases := []struct {
		name       string
		status     core.ContainerStatus
		expID      string
		expRuntime string
	}{
		{
			name: "Docker Container extract runtime",
			status: core.ContainerStatus{
				ContainerID: "docker://fb5dca57a03a05cd7b1291a6cf295196dbfaae51cc5c477ec8748817df4b7208",
			},
			expID:      "fb5dca57a03a05cd7b1291a6cf295196dbfaae51cc5c477ec8748817df4b7208",
			expRuntime: "docker",
		},
		{
			name: "ContainerD Container extract runtime",
			status: core.ContainerStatus{
				ContainerID: "containerd://fb5dca57a03a05cd7b1291a6cf295196dbfaae51cc5c477ec8748817df4b7208",
			},
			expID:      "fb5dca57a03a05cd7b1291a6cf295196dbfaae51cc5c477ec8748817df4b7208",
			expRuntime: "containerd",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			info := NewContainerInfo(&tc.status)

			require.Equal(t, tc.expRuntime, info.Runtime)
			require.Equal(t, tc.expID, info.ID)
		})
	}
}

func Test_GetContainerFSVolume(t *testing.T) {
	testCases := []struct {
		name          string
		status        core.ContainerStatus
		expVolumeName string
		expVolumePath string
	}{
		{
			name: "DockerFS used if specified",
			status: core.ContainerStatus{
				ContainerID: "docker://fb5dca57a03a05cd7b1291a6cf295196dbfaae51cc5c477ec8748817df4b7208",
			},
			expVolumeName: "dockerfs",
			expVolumePath: globals.PathDockerFS,
		},
		{
			name: "ContainerdFS used if specified",
			status: core.ContainerStatus{
				ContainerID: "containerd://fb5dca57a03a05cd7b1291a6cf295196dbfaae51cc5c477ec8748817df4b7208",
			},
			expVolumeName: "containerdfs",
			expVolumePath: globals.PathContainerDFS,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			container := NewContainerInfo(&tc.status)

			volume := container.GetContainerFSVolume()

			require.Equal(t, tc.expVolumeName, volume.Name)
			require.Equal(t, tc.expVolumePath, volume.HostPath)
			require.Equal(t, tc.expVolumePath, volume.MountPath)
		})
	}
}

func Test_GetContainerSharedVolumes(t *testing.T) {
	testCases := []struct {
		name          string
		status        core.ContainerStatus
		expVolumeName string
		expVolumePath string
	}{
		{
			name: "Docker mounts used if specified",
			status: core.ContainerStatus{
				ContainerID: "docker://fb5dca57a03a05cd7b1291a6cf295196dbfaae51cc5c477ec8748817df4b7208",
			},
			expVolumeName: "dockervolumes",
			expVolumePath: globals.PathDockerVolumes,
		},
		{
			name: "Containerd mounts used if specified",
			status: core.ContainerStatus{
				ContainerID: "containerd://fb5dca57a03a05cd7b1291a6cf295196dbfaae51cc5c477ec8748817df4b7208",
			},
			expVolumeName: "containerdvolumes",
			expVolumePath: globals.PathContainerDVolumes,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			container := NewContainerInfo(&tc.status)

			volume := container.GetContainerSharedVolumes()

			require.Equal(t, tc.expVolumeName, volume.Name)
			require.Equal(t, tc.expVolumePath, volume.HostPath)
			require.Equal(t, tc.expVolumePath, volume.MountPath)
		})
	}
}
