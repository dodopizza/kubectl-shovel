package kubernetes

import (
	"github.com/stretchr/testify/require"
	core "k8s.io/api/core/v1"
	"testing"
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
		name       string
		runtime    string
		volumeName string
	}{
		{
			name:       "DockerFS used if specified",
			runtime:    "docker",
			volumeName: "dockerfs",
		},
		{
			name:       "ContainerdFS used if specified",
			runtime:    "containerd",
			volumeName: "containerdfs",
		},
		{
			name:       "DockerFS used in default",
			runtime:    "",
			volumeName: "dockerfs",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			container := &ContainerInfo{Runtime: tc.runtime}

			volume := container.GetContainerFSVolume()

			require.Equal(t, tc.volumeName, volume.Name)
		})
	}
}
