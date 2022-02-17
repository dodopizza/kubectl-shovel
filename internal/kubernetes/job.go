package kubernetes

import (
	"context"
	"github.com/google/uuid"
	"strings"

	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func JobName() string {
	parts := []string{"kubectl-shovel", uuid.NewString()}
	return strings.Join(parts, "-")
}

// RunJob will run job with specified parameters
func (k8s *Client) RunJob(
	name,
	image,
	nodeName string,
	volume *JobVolume,
	cmdArgs []string,
) error {
	imageName := image

	commonMeta := metav1.ObjectMeta{
		Name:      name,
		Namespace: k8s.Namespace,
		Labels: map[string]string{
			"job-name": name,
		},
	}

	job := &batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Job",
			APIVersion: "batch/v1",
		},
		ObjectMeta: commonMeta,
		Spec: batchv1.JobSpec{
			Parallelism:             int32Ptr(1),
			Completions:             int32Ptr(1),
			TTLSecondsAfterFinished: int32Ptr(5),
			BackoffLimit:            int32Ptr(0),
			Template: v1.PodTemplateSpec{
				ObjectMeta: commonMeta,
				Spec: v1.PodSpec{
					Volumes: []apiv1.Volume{
						{
							Name: volume.Name,
							VolumeSource: apiv1.VolumeSource{
								HostPath: &apiv1.HostPathVolumeSource{
									Path: volume.HostPath,
								},
							},
						},
					},
					InitContainers: nil,
					Containers: []apiv1.Container{
						{
							ImagePullPolicy: v1.PullIfNotPresent,
							Name:            "kubectl-shovel",
							Image:           imageName,
							TTY:             true,
							Stdin:           true,
							Args:            cmdArgs,
							VolumeMounts: []apiv1.VolumeMount{
								{
									Name:      volume.Name,
									MountPath: volume.MountPath,
								},
							},
						},
					},
					RestartPolicy: "Never",
					NodeName:      nodeName,
				},
			},
		},
	}

	_, err := k8s.
		BatchV1().
		Jobs(k8s.Namespace).
		Create(context.Background(), job, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}

// DeleteJob deleting job
func (k8s *Client) DeleteJob(name string) error {
	propagationPolicy := metav1.DeletePropagationForeground
	return k8s.
		BatchV1().
		Jobs(k8s.Namespace).
		Delete(
			context.Background(),
			name,
			metav1.DeleteOptions{
				PropagationPolicy: &propagationPolicy,
			},
		)
}
