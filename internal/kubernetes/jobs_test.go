package kubernetes

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	batch "k8s.io/api/batch/v1"
	core "k8s.io/api/core/v1"

	"github.com/dodopizza/kubectl-shovel/internal/globals"
)

var (
	jobName      = fmt.Sprintf("%s-suffix", globals.PluginName)
	jobNamespace = "shovel-namespace"
)

func requireJobSpecMatches(t *testing.T, spec *JobRunSpec, job *batch.Job) {
	require.Equal(t, jobName, job.Name)
	require.Equal(t, jobNamespace, job.Namespace)
	require.Equal(t, spec.Selectors, job.Labels)
	require.Equal(t, spec.Node, job.Spec.Template.Spec.NodeName)
	require.Equal(t, 1, len(job.Spec.Template.Spec.Containers))
	require.Equal(t, globals.PluginName, job.Spec.Template.Spec.Containers[0].Name)
	require.Equal(t, spec.Args, job.Spec.Template.Spec.Containers[0].Args)
	require.Equal(t, spec.Image, job.Spec.Template.Spec.Containers[0].Image)
}

func Test_NewRunJobSpec(t *testing.T) {
	generator = func() string {
		return "suffix"
	}
	testCases := []struct {
		args      []string
		image     string
		name      string
		pod       *PodInfo
		container *ContainerInfo
	}{
		{
			name:  "Correct job spec generates",
			args:  []string{"/bin/sh", "-c", "sleep", "1"},
			image: "alpine",
			pod: NewPodInfo(&core.Pod{
				Spec: core.PodSpec{
					NodeName: "node",
				},
			}),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			spec := NewJobRunSpec(tc.args, tc.image, tc.pod)

			job := spec.Build(jobNamespace)

			requireJobSpecMatches(t, spec, job)
		})
	}
}

func Test_JobRunSpecWithContainerFSVolumes(t *testing.T) {
	testCases := []struct {
		name      string
		container *ContainerInfo
		expCount  int
	}{
		{
			name: "ContainerD volumes added volume mounts per each volume",
			container: NewContainerInfo(
				&core.ContainerStatus{
					ContainerID: "containerd://fb5dca57a03a05cd7b1291a6cf295196dbfaae51cc5c477ec8748817df4b7208",
				},
			),
			expCount: 2,
		},
		{
			name: "Docker volumes added only one",
			container: NewContainerInfo(
				&core.ContainerStatus{
					ContainerID: "docker://fb5dca57a03a05cd7b1291a6cf295196dbfaae51cc5c477ec8748817df4b7208",
				},
			),
			expCount: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pod := NewPodInfo(&core.Pod{
				Spec: core.PodSpec{
					NodeName: "node",
				},
			})

			spec := NewJobRunSpec([]string{"sleep"}, "alpine", pod).
				WithContainerFSVolume(tc.container).
				WithContainerMountsVolume(tc.container)
			job := spec.Build(jobNamespace)

			requireJobSpecMatches(t, spec, job)
			require.Equal(t, tc.expCount, len(job.Spec.Template.Spec.Volumes))
			require.Equal(t, tc.expCount, len(job.Spec.Template.Spec.Containers[0].VolumeMounts))
		})
	}
}

func Test_JobRunSpecWithHostTmpVolume(t *testing.T) {
	pod := NewPodInfo(&core.Pod{
		Spec: core.PodSpec{
			NodeName: "node",
		},
	})

	spec := NewJobRunSpec([]string{"sleep"}, "alpine", pod).
		WithHostTmpVolume("/tmp/testing")
	job := spec.Build(jobNamespace)

	requireJobSpecMatches(t, spec, job)
	require.Equal(t, 1, len(job.Spec.Template.Spec.Volumes))
	require.Equal(t, 1, len(job.Spec.Template.Spec.Containers[0].VolumeMounts))
	require.Equal(t, "hostoutput", job.Spec.Template.Spec.Volumes[0].Name)
	require.Equal(t, "/tmp/testing", job.Spec.Template.Spec.Volumes[0].HostPath.Path)
	require.Equal(t, globals.PathHostOutputFolder, job.Spec.Template.Spec.Containers[0].VolumeMounts[0].MountPath)
}

func Test_JobRunSpecWithPrivileged(t *testing.T) {
	pod := NewPodInfo(&core.Pod{
		Spec: core.PodSpec{
			NodeName: "node",
		},
	})

	spec := NewJobRunSpec([]string{"sleep"}, "alpine", pod).
		WithHostProcVolume().
		WithPrivilegedOptions()
	job := spec.Build(jobNamespace)

	requireJobSpecMatches(t, spec, job)
	require.Equal(t, 1, len(job.Spec.Template.Spec.Volumes))
	require.Equal(t, 1, len(job.Spec.Template.Spec.Containers[0].VolumeMounts))
	require.Equal(t, "hostproc", job.Spec.Template.Spec.Volumes[0].Name)
	require.Equal(t, globals.PathHostProcFolder, job.Spec.Template.Spec.Volumes[0].HostPath.Path)
	require.Equal(t, globals.PathHostProcFolder, job.Spec.Template.Spec.Containers[0].VolumeMounts[0].MountPath)
	require.NotEqual(t, nil, job.Spec.Template.Spec.Containers[0].SecurityContext)
	require.Equal(t, true, *job.Spec.Template.Spec.Containers[0].SecurityContext.Privileged)
	require.Equal(t,
		[]core.Capability{"SYS_PTRACE"},
		job.Spec.Template.Spec.Containers[0].SecurityContext.Capabilities.Add)
	require.Equal(t, true, job.Spec.Template.Spec.HostPID)
}
