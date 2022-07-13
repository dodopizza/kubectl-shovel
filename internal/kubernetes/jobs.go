package kubernetes

import (
	"context"
	"strings"

	"github.com/google/uuid"

	"github.com/dodopizza/kubectl-shovel/internal/globals"

	batch "k8s.io/api/batch/v1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	// add ability to test behavior
	generator = uuid.NewString
)

// JobRunSpec is helper struct to describe job spec customization
type JobRunSpec struct {
	Args       []string
	Image      string
	Name       string
	Node       string
	Privileged bool
	Selectors  map[string]string
	Volumes    []JobVolume
}

// JobVolume is helper struct to describe job volume customization
type JobVolume struct {
	Name      string
	HostPath  string
	MountPath string
}

// NewJobRunSpec returns JobRunSpec constructed from specified args, image, pod and container
func NewJobRunSpec(args []string, image string, pod *PodInfo) *JobRunSpec {
	nameParts := []string{globals.PluginName, generator()}
	name := strings.Join(nameParts, "-")

	return &JobRunSpec{
		Args:  args,
		Image: image,
		Name:  name,
		Node:  pod.Node,
		Selectors: map[string]string{
			"job-name": name,
		},
		Volumes: []JobVolume{},
	}
}

// WithPrivilegedOptions adds core.SecurityContext to with SYS_PTRACE capability to core.JobSpec
func (j *JobRunSpec) WithPrivilegedOptions() *JobRunSpec {
	j.Privileged = true
	return j
}

// WithContainerFSVolume add host volume that used to store container file system volumes
func (j *JobRunSpec) WithContainerFSVolume(container *ContainerInfo) *JobRunSpec {
	j.appendVolume(container.GetContainerFSVolume())
	return j
}

// WithContainerMountsVolume add host volume that used to store container additional volumes
func (j *JobRunSpec) WithContainerMountsVolume(container *ContainerInfo) *JobRunSpec {
	j.appendVolume(container.GetContainerSharedVolumes())
	return j
}

// WithHostProcVolume adds host volume that used to locate target process memory sections
func (j *JobRunSpec) WithHostProcVolume() *JobRunSpec {
	j.appendVolume(JobVolume{
		Name:      "hostproc",
		HostPath:  globals.PathHostProcFolder,
		MountPath: globals.PathHostProcFolder,
	})
	return j
}

// WithHostTmpVolume add host /tmp volume that used to store output
func (j *JobRunSpec) WithHostTmpVolume(path string) *JobRunSpec {
	j.appendVolume(JobVolume{
		Name:      "hostoutput",
		HostPath:  path,
		MountPath: globals.PathHostOutputFolder,
	})
	return j
}

// Build returns resulted batch.Job spec from JobRunSpec options
func (j *JobRunSpec) Build(namespace string) *batch.Job {
	metaSpec := meta.ObjectMeta{
		Name:      j.Name,
		Namespace: namespace,
		Labels:    j.Selectors,
	}

	return &batch.Job{
		TypeMeta: meta.TypeMeta{
			Kind:       "Job",
			APIVersion: "batch/v1",
		},
		ObjectMeta: metaSpec,
		Spec: batch.JobSpec{
			Parallelism:             ptr(int32(1)),
			Completions:             ptr(int32(1)),
			TTLSecondsAfterFinished: ptr(int32(5)),
			BackoffLimit:            ptr(int32(0)),
			Template: core.PodTemplateSpec{
				ObjectMeta: metaSpec,
				Spec: core.PodSpec{
					Volumes:        j.volumes(),
					InitContainers: nil,
					Containers: []core.Container{
						{
							ImagePullPolicy:          core.PullIfNotPresent,
							Name:                     globals.PluginName,
							Image:                    j.Image,
							TTY:                      true,
							Stdin:                    true,
							Args:                     j.Args,
							VolumeMounts:             j.mounts(),
							TerminationMessagePolicy: core.TerminationMessageFallbackToLogsOnError,
							SecurityContext:          j.securityContext(),
						},
					},
					HostPID:       j.Privileged,
					RestartPolicy: "Never",
					NodeName:      j.Node,
				},
			},
		},
	}
}

func (j *JobRunSpec) appendVolume(item JobVolume) {
	// ignore any duplicates by host path
	for _, volume := range j.Volumes {
		if volume.HostPath == item.HostPath {
			return
		}
	}

	j.Volumes = append(j.Volumes, item)
}

func (j *JobRunSpec) volumes() []core.Volume {
	volumes := make([]core.Volume, len(j.Volumes))
	for i, volume := range j.Volumes {
		volumes[i] = core.Volume{
			Name: volume.Name,
			VolumeSource: core.VolumeSource{
				HostPath: &core.HostPathVolumeSource{
					Path: volume.HostPath,
				},
			},
		}
	}
	return volumes
}

func (j *JobRunSpec) mounts() []core.VolumeMount {
	volumeMounts := make([]core.VolumeMount, len(j.Volumes))
	for i, volume := range j.Volumes {
		volumeMounts[i] = core.VolumeMount{
			Name:      volume.Name,
			MountPath: volume.MountPath,
		}
	}
	return volumeMounts
}

func (j *JobRunSpec) securityContext() *core.SecurityContext {
	if !j.Privileged {
		return nil
	}

	return &core.SecurityContext{
		Capabilities: &core.Capabilities{
			Add: []core.Capability{"SYS_PTRACE"},
		},
		Privileged: ptr(true),
	}
}

// RunJob will run job with specified parameters
func (k *Client) RunJob(spec *JobRunSpec) error {
	job := spec.Build(k.Namespace)

	_, err := k.
		BatchV1().
		Jobs(k.Namespace).
		Create(context.Background(), job, meta.CreateOptions{})

	return err
}

// DeleteJob deleting job
func (k *Client) DeleteJob(name string) error {
	policy := meta.DeletePropagationForeground
	options := meta.DeleteOptions{
		PropagationPolicy: &policy,
	}

	return k.
		BatchV1().
		Jobs(k.Namespace).
		Delete(context.Background(), name, options)
}

func ptr[T any](t T) *T {
	return &t
}
