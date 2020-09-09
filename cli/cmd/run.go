package cmd

import (
	"fmt"
	"strings"

	"github.com/dodopizza/kubectl-shovel/cli/kubernetes"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func run(
	kubeFlags *genericclioptions.ConfigFlags,
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
	containerID := strings.TrimPrefix(pod.Status.ContainerStatuses[0].ContainerID, "docker://")
	fmt.Println("Run diagnostics job")
	err = k8s.RunJob(
		jobName,
		dumperImageName,
		pod.Spec.NodeName,
		[]string{
			"--container-id",
			containerID,
			"--tool",
			tool,
		},
	)
	if err != nil {
		return err
	}

	fmt.Println("Waiting job to start")
	jobPodName, err := k8s.WaitPod(jobName)
	if err != nil {
		return err
	}
	readCloser, err := k8s.ReadPodLogs(jobPodName)
	if err != nil {
		return err
	}

	handleLogs(readCloser, output)
	fmt.Printf("Result successfuly written to %s\n", output)

	err = k8s.DeleteJob(jobName)

	return nil
}
