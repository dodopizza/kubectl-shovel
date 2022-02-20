package kubernetes

import (
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	"testing"
)

func Test_NewContainerInfo(t *testing.T) {
	testCases := []struct {
		name       string
		status     v1.ContainerStatus
		expID      string
		expRuntime string
	}{
		{
			name: "Docker Container extract runtime",
			status: v1.ContainerStatus{
				ContainerID: "docker://fb5dca57a03a05cd7b1291a6cf295196dbfaae51cc5c477ec8748817df4b7208",
			},
			expID:      "fb5dca57a03a05cd7b1291a6cf295196dbfaae51cc5c477ec8748817df4b7208",
			expRuntime: "docker",
		},
		{
			name: "ContainerD Container extract runtime",
			status: v1.ContainerStatus{
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
