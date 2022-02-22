package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/dodopizza/kubectl-shovel/internal/events"
	"github.com/dodopizza/kubectl-shovel/internal/globals"

	"github.com/pkg/errors"

	"github.com/dodopizza/kubectl-shovel/internal/kubernetes"
	"github.com/dodopizza/kubectl-shovel/internal/watchdog"
)

func (cb *CommandBuilder) newKubeClient() error {
	kube, err := kubernetes.NewClient(cb.CommonOptions.kubeConfig)
	if err != nil {
		return err
	}

	cb.kube = kube
	return nil
}

func (cb *CommandBuilder) args(pod *kubernetes.PodInfo, container *kubernetes.ContainerInfo) []string {
	args := []string{"--container-id", container.ID, "--container-runtime", container.Runtime}

	if cb.CommonOptions.StoreOutputOnHost {
		args = append(args, "store-output-on-host")
	}

	args = append(args, "--container-name", container.Name)
	args = append(args, "--pod-name", pod.Name)
	args = append(args, "--pod-namespace", pod.Namespace)
	args = append(args, cb.tool.ToolName())
	args = append(args, cb.tool.FormatArgs()...)
	return args
}

func (cb *CommandBuilder) copyOutputFromJob(pod *kubernetes.PodInfo, output string) error {
	fmt.Println("Retrieve output from diagnostics job")
	if err := cb.kube.CopyFromPod(pod.Name, output, cb.CommonOptions.Output); err != nil {
		return errors.Wrap(err, "Error while retrieving diagnostics job output")
	}
	fmt.Printf("Result successfully written to %s\n", cb.CommonOptions.Output)
	return nil
}

func (cb *CommandBuilder) storeOutputOnHost(pod *kubernetes.PodInfo, output string) error {
	fmt.Printf("Output located on host: %s, at path: %s\n", pod.Node, output)
	return nil
}

func (cb *CommandBuilder) launch() error {
	if err := cb.newKubeClient(); err != nil {
		return errors.Wrap(err, "Failed to init kubernetes client")
	}

	targetPod, err := cb.kube.GetPodInfo(cb.CommonOptions.Pod)
	if err != nil {
		return errors.Wrap(err, "Failed to get info about target pod")
	}

	targetContainerName := cb.CommonOptions.Container
	if targetContainerName == "" {
		targetContainerName = targetPod.Annotations["kubectl.kubernetes.io/default-container"]
	}

	targetContainer, err := targetPod.FindContainerInfo(targetContainerName)
	if err != nil {
		return errors.Wrap(err, "Failed to get info about target container")
	}

	jobSpec := kubernetes.
		NewJobRunSpec(cb.args(targetPod, targetContainer), cb.CommonOptions.Image, targetPod).
		WithContainerFSVolume(targetContainer)

	if targetPod.ContainsMountedTmp(targetContainerName) {
		jobSpec.WithContainerMountsVolume(targetContainer)
	}

	if cb.CommonOptions.StoreOutputOnHost {
		jobSpec.WithHostTmpVolume()
	}

	fmt.Printf("Spawning diagnostics job with command:\n%s\n", strings.Join(jobSpec.Args, " "))
	if err := cb.kube.RunJob(jobSpec); err != nil {
		return errors.Wrap(err, "Failed to spawn diagnostics job")
	}

	fmt.Println("Waiting for a diagnostics job to start")
	jobPod, err := cb.kube.WaitPod(jobSpec.Selectors)
	if err != nil {
		return errors.Wrap(err, "Failed to wait diagnostics job execution")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pinger := watchdog.NewPinger(cb.kube, jobPod.Name)
	go func() {
		if err := pinger.Run(ctx); err != nil {
			fmt.Println(err)
		}
	}()

	jobPodLogs, err := cb.kube.ReadPodLogs(jobPod.Name, globals.PluginName)
	if err != nil {
		return errors.Wrap(err, "Failed to read logs from diagnostics job targetPod")
	}
	defer jobPodLogs.Close()

	awaiter := events.NewEventAwaiter()
	output, err := awaiter.AwaitCompletedEvent(jobPodLogs)
	if err != nil {
		return err
	}

	// dealing with output
	outputHandler := cb.copyOutputFromJob
	if cb.CommonOptions.StoreOutputOnHost {
		outputHandler = cb.storeOutputOnHost
	}
	if err := outputHandler(jobPod, output); err != nil {
		return err
	}

	fmt.Println("Cleanup diagnostics job")
	if err := cb.kube.DeleteJob(jobSpec.Name); err != nil {
		return errors.Wrap(err, "Error while deleting job")
	}

	return nil
}
