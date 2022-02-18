package cmd

import (
	"context"
	"fmt"
	"github.com/dodopizza/kubectl-shovel/internal/events"
	"github.com/dodopizza/kubectl-shovel/internal/globals"

	"github.com/pkg/errors"

	"github.com/dodopizza/kubectl-shovel/internal/kubernetes"
	"github.com/dodopizza/kubectl-shovel/internal/watchdog"
)

func (cb *CommandBuilder) launch() error {
	k8s, err := kubernetes.NewClient(cb.CommonOptions.kubeFlags)
	if err != nil {
		return nil
	}

	pod, err := k8s.GetPodInfo(cb.CommonOptions.podName)
	if err != nil {
		return err
	}

	containerInfo, err := kubernetes.GetContainerInfo(pod, cb.CommonOptions.containerName)
	if err != nil {
		return errors.Wrap(err, "Error while getting info about container")
	}

	jobName := kubernetes.JobName()
	jobVolume := containerInfo.GetJobVolume()

	fmt.Println("Spawning diagnostics job")
	args := append([]string{
		"--container-id",
		containerInfo.ID,
		"--container-runtime",
		containerInfo.Runtime,
		cb.tool.ToolName(),
	}, cb.tool.GetArgs()...)
	if err := k8s.RunJob(
		jobName,
		cb.CommonOptions.image,
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

	stream, err := k8s.ReadPodLogs(jobPodName, globals.PluginName)
	if err != nil {
		return err
	}
	defer stream.Close()

	awaiter := events.NewEventAwaiter()
	resultFilePath, err := awaiter.AwaitCompletedEvent(stream)
	if err != nil {
		return err
	}

	fmt.Println("Getting results from job")
	if err := k8s.CopyFromPod(jobPodName, resultFilePath, cb.CommonOptions.output); err != nil {
		return errors.Wrap(err, "Error while getting results")
	}
	fmt.Printf("Result successfully written to %s\n", cb.CommonOptions.output)

	if err := k8s.DeleteJob(jobName); err != nil {
		return errors.Wrap(err, "Error while deleting job")
	}

	return nil
}
