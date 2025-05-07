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
		containers    []core.Container
		podStatus     core.PodStatus
		containerName string
		expRuntime    string
		expID         string
	}{
		{
			name:       "Docker container",
			containers: []core.Container{{Name: "target"}},
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
			name:       "Containerd container",
			containers: []core.Container{{Name: "target"}},
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
			name:       "Specified container name",
			containers: []core.Container{{Name: "target"}},
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
			name:       "MultiContainer pod",
			containers: []core.Container{{Name: "target"}, {Name: "wrong"}},
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
			name:       "MultiContainer pod",
			containers: []core.Container{{Name: "wrong"}, {Name: "target"}},
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
				Spec: core.PodSpec{
					Containers: tc.containers,
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

func Test_GetContainerInfo_InitContainer(t *testing.T) {
	testCases := []struct {
		name          string
		containers    []core.Container
		initContainers []core.Container
		podStatus     core.PodStatus
		containerName string
		expRuntime    string
		expID         string
		expIsInit     bool
	}{
		{
			name:       "Find init container when specified directly",
			containers: []core.Container{{Name: "app"}},
			initContainers: []core.Container{{Name: "init-container"}},
			podStatus: core.PodStatus{
				ContainerStatuses: []core.ContainerStatus{
					{
						Name:        "app",
						ContainerID: "docker://regular-container-id",
					},
				},
				InitContainerStatuses: []core.ContainerStatus{
					{
						Name:        "init-container",
						ContainerID: "docker://init-container-id",
					},
				},
			},
			containerName: "init-container",
			expRuntime:    "docker",
			expID:         "init-container-id",
			expIsInit:     true,
		},
		{
			name:       "Fallback to init container when not found in regular containers",
			containers: []core.Container{{Name: "app"}},
			initContainers: []core.Container{{Name: "side-container"}},
			podStatus: core.PodStatus{
				ContainerStatuses: []core.ContainerStatus{
					{
						Name:        "app",
						ContainerID: "docker://regular-container-id",
					},
				},
				InitContainerStatuses: []core.ContainerStatus{
					{
						Name:        "side-container",
						ContainerID: "docker://side-container-id",
					},
				},
			},
			containerName: "side-container",
			expRuntime:    "docker",
			expID:         "side-container-id",
			expIsInit:     true,
		},
		{
			name:       "Regular container takes precedence over init container with same name",
			containers: []core.Container{{Name: "same-name"}},
			initContainers: []core.Container{{Name: "same-name"}},
			podStatus: core.PodStatus{
				ContainerStatuses: []core.ContainerStatus{
					{
						Name:        "same-name",
						ContainerID: "docker://regular-container-id",
					},
				},
				InitContainerStatuses: []core.ContainerStatus{
					{
						Name:        "same-name",
						ContainerID: "docker://init-container-id",
					},
				},
			},
			containerName: "same-name",
			expRuntime:    "docker",
			expID:         "regular-container-id",
			expIsInit:     false,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			podInfo := NewPodInfo(&core.Pod{
				ObjectMeta: meta.ObjectMeta{
					Name: tc.name,
				},
				Spec: core.PodSpec{
					Containers: tc.containers,
					InitContainers: tc.initContainers,
				},
				Status: tc.podStatus,
			})

			info, err := podInfo.FindContainerInfo(tc.containerName)

			require.NoError(t, err)
			require.Equal(t, tc.expRuntime, info.Runtime)
			require.Equal(t, tc.expID, info.ID)
			require.Equal(t, tc.expIsInit, podInfo.IsInitContainer(tc.containerName))
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
