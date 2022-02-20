package kubernetes

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_GetContainerJobVolume(t *testing.T) {
	testCases := []struct {
		runtime    string
		volumeName string
	}{
		{
			runtime:    "docker",
			volumeName: "dockerfs",
		},
		{
			runtime:    "",
			volumeName: "dockerfs",
		},
		{
			runtime:    "containerd",
			volumeName: "containerdfs",
		},
	}

	for _, tc := range testCases {
		container := &ContainerInfo{Runtime: tc.runtime}

		volumes := container.GetJobVolumes()
		volume := volumes[0]

		require.Equal(t, tc.volumeName, volume.Name)
	}
}
