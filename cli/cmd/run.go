package cmd

import (
	"fmt"

	"github.com/dodopizza/kubectl-shovel/internal/kubernetes"
	"github.com/dodopizza/kubectl-shovel/internal/watchdog"
	"github.com/pkg/errors"
)

func run(
	options *commonOptions,
	tool string,
) error {
	k8s, err := kubernetes.NewClient(options.kubeFlags)
	if err != nil {
		return nil
	}
	pod, err := k8s.GetPodInfo(options.podName)
	if err != nil {
		return err
	}

	jobName := newJobName()
	containerInfo := kubernetes.GetContainerInfo(pod)
	jobVolume := kubernetes.NewJobVolume(containerInfo)
	fmt.Println("Spawning diagnostics job")
	err = k8s.RunJob(
		jobName,
		options.image,
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
	jobPodName, err := k8s.WaitPod(map[string]string{"job-name": jobName})
	if err != nil {
		return err
	}

	op := watchdog.NewOperator(k8s, jobPodName)
	go func() {
		if err := op.Run(); err != nil {
			fmt.Println(err)
		}
	}()

	readCloser, err := k8s.ReadPodLogs(jobPodName)
	if err != nil {
		return err
	}

	var resultFilePath string
	if resultFilePath, err = handleLogs(readCloser); err != nil {
		return err
	}

	fmt.Println("Getting results from job")
	if err := k8s.CopyFromPod(jobPodName, resultFilePath, options.output); err != nil {
		return errors.Wrap(err, "Error while getting results")
	}
	fmt.Printf("Result successfully written to %s\n", options.output)

	if err := k8s.DeleteJob(jobName); err != nil {
		return errors.Wrap(err, "Error while deleting job")
	}

	return nil
}
