package cmd

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/dodopizza/kubectl-shovel/internal/kubernetes"
	"github.com/dodopizza/kubectl-shovel/internal/watchdog"
)

func run(
	options *commonOptions,
	tool string,
	args ...string,
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
	containerInfo, err := kubernetes.GetContainerInfo(
		pod,
		options.containerName,
	)
	if err != nil {
		return errors.Wrap(err, "Error while getting info about container")
	}

	jobVolume := kubernetes.NewJobVolume(containerInfo)

	fmt.Println("Spawning diagnostics job")
	args = append([]string{
		tool,
		"--container-id",
		containerInfo.ID,
		"--container-runtime",
		containerInfo.Runtime,
	}, args...)
	if err := k8s.RunJob(
		jobName,
		options.image,
		pod.Spec.NodeName,
		jobVolume,
		args,
	); err != nil {
		return err
	}

	fmt.Println("Waiting for a job to start")
	jobPodName, err := k8s.WaitPod(map[string]string{"job-name": jobName})
	if err != nil {
		return err
	}

	op := watchdog.NewOperator(k8s, jobPodName)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		if err := op.Run(ctx); err != nil {
			fmt.Println(err)
		}
	}()

	readCloser, err := k8s.ReadPodLogs(jobPodName, "kubectl-shovel")
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
