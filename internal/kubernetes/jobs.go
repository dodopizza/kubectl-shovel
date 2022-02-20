package kubernetes

import (
	"context"
	"github.com/dodopizza/kubectl-shovel/internal/globals"
	"github.com/google/uuid"
	"strings"

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
	Args      []string
	Image     string
	Name      string
	Node      string
	Selectors map[string]string
	Volumes   []JobVolume
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

// WithContainerFSVolume add container file system volumes to job spec
func (j *JobRunSpec) WithContainerFSVolume(container *ContainerInfo) *JobRunSpec {
	j.Volumes = append(j.Volumes, container.GetContainerFSVolume())
	return j
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

// RunJob will run job with specified parameters
func (k8s *Client) RunJob(spec *JobRunSpec) error {
	commonMeta := meta.ObjectMeta{
		Name:      spec.Name,
		Namespace: k8s.Namespace,
		Labels:    spec.Selectors,
	}

	job := &batch.Job{
		TypeMeta: meta.TypeMeta{
			Kind:       "Job",
			APIVersion: "batch/v1",
		},
		ObjectMeta: commonMeta,
		Spec: batch.JobSpec{
			Parallelism:             int32Ptr(1),
			Completions:             int32Ptr(1),
			TTLSecondsAfterFinished: int32Ptr(5),
			BackoffLimit:            int32Ptr(0),
			Template: core.PodTemplateSpec{
				ObjectMeta: commonMeta,
				Spec: core.PodSpec{
					Volumes:        spec.volumes(),
					InitContainers: nil,
					Containers: []core.Container{
						{
							ImagePullPolicy: core.PullIfNotPresent,
							Name:            globals.PluginName,
							Image:           spec.Image,
							TTY:             true,
							Stdin:           true,
							Args:            spec.Args,
							VolumeMounts:    spec.mounts(),
						},
					},
					RestartPolicy: "Never",
					NodeName:      spec.Node,
				},
			},
		},
	}

	_, err := k8s.
		BatchV1().
		Jobs(k8s.Namespace).
		Create(context.Background(), job, meta.CreateOptions{})

	return err
}

// DeleteJob deleting job
func (k8s *Client) DeleteJob(name string) error {
	propagationPolicy := meta.DeletePropagationForeground
	return k8s.
		BatchV1().
		Jobs(k8s.Namespace).
		Delete(
			context.Background(),
			name,
			meta.DeleteOptions{
				PropagationPolicy: &propagationPolicy,
			},
		)
}

func int32Ptr(i int32) *int32 {
	return &i
}
