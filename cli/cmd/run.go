package cmd

import (
	"fmt"

	"github.com/dodopizza/kubectl-shovel/internal/kubernetes"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func run(
	kubeFlags *genericclioptions.ConfigFlags,
	image,
	podName,
	output,
	tool string,
) error {
	k8s, err := kubernetes.NewClient(kubeFlags)
	if err != nil {
		return nil
	}
	pod, err := k8s.GetPodInfo(podName)
	if err != nil {
		return err
	}

	jobName := newJobName()
	containerInfo := kubernetes.GetContainerInfo(pod)
	jobVolume := kubernetes.NewJobVolume(containerInfo)
	fmt.Println("Spawning diagnostics job")
	err = k8s.RunJob(
		jobName,
		image,
		pod.Spec.NodeName,
		jobVolume,
		[]string{
			tool,
			"--container-id",
			containerInfo.ID,
			"--container-runtime",
			containerInfo.Runtime,
		},
	)
	if err != nil {
		return err
	}

	fmt.Println("Waiting for a job to start")
	jobPodName, err := k8s.WaitPod(jobName)
	if err != nil {
		return err
	}
	readCloser, err := k8s.ReadPodLogs(jobPodName)
	if err != nil {
		return err
	}

	handleLogs(readCloser, output)
	fmt.Printf("Result successfully written to %s\n", output)

	err = k8s.DeleteJob(jobName)
	if err != nil {
		return err
	}

	return nil
}
