package kubernetes

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
)

func Test_GetContainerInfo(t *testing.T) {
	testCases := []ContainerInfo{
		{
			Runtime: "docker",
			ID:      "fb5dca57a03a05cd7b1291a6cf295196dbfaae51cc5c477ec8748817df4b7208",
		},
		{
			Runtime: "containerd",
			ID:      "2202fc17c16fb85a3bba5395278b8b5478154f023981be57edb82d931472f4ac",
		},
	}
	for _, tc := range testCases {
		info := GetContainerInfo(&v1.Pod{
			Status: v1.PodStatus{
				ContainerStatuses: []v1.ContainerStatus{
					{
						ContainerID: fmt.Sprintf("%s://%s", tc.Runtime, tc.ID),
					},
				},
			},
		})

		require.Equal(t, tc.Runtime, info.Runtime)
		require.Equal(t, tc.ID, info.ID)
	}
}

func Test_NewJobVolume(t *testing.T) {
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
		volume := NewJobVolume(&ContainerInfo{
			Runtime: tc.runtime,
		})
		require.Equal(t, tc.volumeName, volume.Name)
	}
}
