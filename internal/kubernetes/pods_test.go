package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/require"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_GetContainerInfo(t *testing.T) {
	testCases := []struct {
		name          string
		podStatus     core.PodStatus
		containerName string
		expRuntime    string
		expID         string
	}{
		{
			name: "Docker container",
			podStatus: core.PodStatus{
				ContainerStatuses: []core.ContainerStatus{
					{
						Name:        "target",
						ContainerID: "docker://fb5dca57a03a05cd7b1291a6cf295196dbfaae51cc5c477ec8748817df4b7208",
					},
				},
			},
			expRuntime: "docker",
			expID:      "fb5dca57a03a05cd7b1291a6cf295196dbfaae51cc5c477ec8748817df4b7208",
		},
		{
			name: "Containerd container",
			podStatus: core.PodStatus{
				ContainerStatuses: []core.ContainerStatus{
					{
						Name:        "target",
						ContainerID: "containerd://2202fc17c16fb85a3bba5395278b8b5478154f023981be57edb82d931472f4ac",
					},
				},
			},
			expRuntime: "containerd",
			expID:      "2202fc17c16fb85a3bba5395278b8b5478154f023981be57edb82d931472f4ac",
		},
		{
			name: "Specified container name",
			podStatus: core.PodStatus{
				ContainerStatuses: []core.ContainerStatus{
					{
						Name:        "target",
						ContainerID: "containerd://2202fc17c16fb85a3bba5395278b8b5478154f023981be57edb82d931472f4ac",
					},
				},
			},
			containerName: "target",
			expRuntime:    "containerd",
			expID:         "2202fc17c16fb85a3bba5395278b8b5478154f023981be57edb82d931472f4ac",
		},
		{
			name: "Multicontainer pod",
			podStatus: core.PodStatus{
				ContainerStatuses: []core.ContainerStatus{
					{
						Name:        "target",
						ContainerID: "containerd://fb5dca57a03a05cd7b1291a6cf295196dbfaae51cc5c477ec8748817df4b7208",
					},
					{
						Name:        "wrong",
						ContainerID: "containerd://2202fc17c16fb85a3bba5395278b8b5478154f023981be57edb82d931472f4ac",
					},
				},
			},
			containerName: "target",
			expRuntime:    "containerd",
			expID:         "fb5dca57a03a05cd7b1291a6cf295196dbfaae51cc5c477ec8748817df4b7208",
		},
		{
			name: "Multicontainer pod",
			podStatus: core.PodStatus{
				ContainerStatuses: []core.ContainerStatus{
					{
						Name:        "wrong",
						ContainerID: "containerd://fb5dca57a03a05cd7b1291a6cf295196dbfaae51cc5c477ec8748817df4b7208",
					},
					{
						Name:        "target",
						ContainerID: "containerd://2202fc17c16fb85a3bba5395278b8b5478154f023981be57edb82d931472f4ac",
					},
				},
			},
			containerName: "target",
			expRuntime:    "containerd",
			expID:         "2202fc17c16fb85a3bba5395278b8b5478154f023981be57edb82d931472f4ac",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			podInfo := NewPodInfo(&core.Pod{
				ObjectMeta: meta.ObjectMeta{
					Name: tc.name,
				},
				Status: tc.podStatus,
			})

			info, err := podInfo.FindContainerInfo(tc.containerName)

			require.NoError(t, err)
			require.Equal(t, tc.expRuntime, info.Runtime)
			require.Equal(t, tc.expID, info.ID)
		})
	}
}

func Test_GetContainerInfo_Error(t *testing.T) {
	testCases := []struct {
		name          string
		podStatus     core.PodStatus
		containerName string
	}{
		{
			name: "Wrong container name",
			podStatus: core.PodStatus{
				ContainerStatuses: []core.ContainerStatus{
					{
						Name:        "target",
						ContainerID: "docker://fb5dca57a03a05cd7b1291a6cf295196dbfaae51cc5c477ec8748817df4b7208",
					},
				},
			},
			containerName: "wrong",
		},
		{
			name: "Empty container name for multicontainer pod",
			podStatus: core.PodStatus{
				ContainerStatuses: []core.ContainerStatus{
					{
						Name:        "first",
						ContainerID: "docker://fb5dca57a03a05cd7b1291a6cf295196dbfaae51cc5c477ec8748817df4b7208",
					},
					{
						Name:        "second",
						ContainerID: "docker://2202fc17c16fb85a3bba5395278b8b5478154f023981be57edb82d931472f4ac",
					},
				},
			},
			containerName: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			podInfo := NewPodInfo(&core.Pod{
				ObjectMeta: meta.ObjectMeta{
					Name: tc.name,
				},
				Status: tc.podStatus,
			})

			_, err := podInfo.FindContainerInfo(tc.containerName)

			require.Error(t, err)
		})
	}
}
