package cmd

import (
	"context"
	"fmt"
	"github.com/dodopizza/kubectl-shovel/internal/events"
	"github.com/dodopizza/kubectl-shovel/internal/globals"
	"strings"

	"github.com/pkg/errors"

	"github.com/dodopizza/kubectl-shovel/internal/kubernetes"
	"github.com/dodopizza/kubectl-shovel/internal/watchdog"
)

func (cb *CommandBuilder) args(info *kubernetes.ContainerInfo) []string {
	args := []string{"--container-id", info.ID, "--container-runtime", info.Runtime}
	args = append(args, cb.tool.ToolName())
	args = append(args, cb.tool.FormatArgs()...)
	return args
}

func (cb *CommandBuilder) launch() error {
	k8s, err := kubernetes.NewClient(cb.CommonOptions.kube)
	if err != nil {
		return errors.Wrap(err, "Failed to init kubernetes client")
	}

	targetPod, err := k8s.GetPodInfo(cb.CommonOptions.Pod)
	if err != nil {
		return errors.Wrap(err, "Failed to get info about target pod")
	}

	targetContainer, err := targetPod.FindContainerInfo(cb.CommonOptions.Container)
	if err != nil {
		return errors.Wrap(err, "Failed to get info about target container")
	}

	jobSpec := kubernetes.
		NewJobRunSpec(cb.args(targetContainer), cb.CommonOptions.Image, targetPod).
		WithContainerFSVolume(targetContainer)

	if targetPod.ContainsMountedTmp(cb.CommonOptions.Container) {
		jobSpec.WithContainerMountsVolume(targetContainer)
	}

	fmt.Printf("Spawning diagnostics job with command:\n%s\n", strings.Join(jobSpec.Args, " "))
	if err := k8s.RunJob(jobSpec); err != nil {
		return errors.Wrap(err, "Failed to spawn diagnostics job")
	}

	fmt.Println("Waiting for a diagnostics job to start")
	jobPod, err := k8s.WaitPod(jobSpec.Selectors)
	if err != nil {
		return errors.Wrap(err, "Failed to wait diagnostics job execution")
	}

	op := watchdog.NewOperator(k8s, jobPod.Name)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		if err := op.Run(ctx); err != nil {
			fmt.Println(err)
		}
	}()

	jobPodLogs, err := k8s.ReadPodLogs(jobPod.Name, globals.PluginName)
	if err != nil {
		return errors.Wrap(err, "Failed to read logs from diagnostics job targetPod")
	}
	defer jobPodLogs.Close()

	awaiter := events.NewEventAwaiter()
	output, err := awaiter.AwaitCompletedEvent(jobPodLogs)
	if err != nil {
		return err
	}

	fmt.Println("Retrieve output from diagnostics job")
	if err := k8s.CopyFromPod(jobPod.Name, output, cb.CommonOptions.Output); err != nil {
		return errors.Wrap(err, "Error while retrieving diagnostics job output")
	}

	fmt.Printf("Result successfully written to %s\nCleanup diagnostics job", cb.CommonOptions.Output)
	if err := k8s.DeleteJob(jobSpec.Name); err != nil {
		return errors.Wrap(err, "Error while deleting job")
	}

	return nil
}
