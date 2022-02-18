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
	k8s, err := kubernetes.NewClient(cb.CommonOptions.kube)
	if err != nil {
		return errors.Wrap(err, "Failed to init kubernetes client")
	}

	pod, err := k8s.GetPodInfo(cb.CommonOptions.Pod)
	if err != nil {
		return errors.Wrap(err, "Failed to get info about target pod")
	}

	containerInfo, err := kubernetes.GetContainerInfo(pod, cb.CommonOptions.Container)
	if err != nil {
		return errors.Wrap(err, "Failed to get info about container")
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
		cb.CommonOptions.Image,
		pod.Spec.NodeName,
		jobVolume,
		args,
	); err != nil {
		return errors.Wrap(err, "Failed to spawn diagnostics job")
	}

	fmt.Println("Waiting for a diagnostics job to start")
	jobPodName, err := k8s.WaitPod(map[string]string{"job-name": jobName})
	if err != nil {
		return errors.Wrap(err, "Failed to wait diagnostics job execution")
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
		return errors.Wrap(err, "Failed to read logs from diagnostics job pod")
	}
	defer stream.Close()

	awaiter := events.NewEventAwaiter()
	resultFilePath, err := awaiter.AwaitCompletedEvent(stream)
	if err != nil {
		return err
	}

	fmt.Println("Retrieve output from diagnostics job")
	if err := k8s.CopyFromPod(jobPodName, resultFilePath, cb.CommonOptions.Output); err != nil {
		return errors.Wrap(err, "Error while retrieving diagnostics job output")
	}

	fmt.Printf("Result successfully written to %s\nCleanup diagnostics job", cb.CommonOptions.Output)
	if err := k8s.DeleteJob(jobName); err != nil {
		return errors.Wrap(err, "Error while deleting job")
	}

	return nil
}
