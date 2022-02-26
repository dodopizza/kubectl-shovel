package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/dodopizza/kubectl-shovel/internal/events"
	"github.com/dodopizza/kubectl-shovel/internal/flags"
	"github.com/dodopizza/kubectl-shovel/internal/globals"
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
	args := flags.NewArgs().
		Append("container-id", container.ID).
		Append("container-runtime", container.Runtime)

	if cb.CommonOptions.StoreOutputOnHost {
		args.AppendKey("store-output-on-host")
	}

	return args.
		Append("container-name", container.Name).
		Append("pod-name", pod.Name).
		Append("pod-namespace", pod.Namespace).
		AppendCommand(cb.tool.ToolName()).
		AppendFrom(cb.tool).
		Get()
}

func (cb *CommandBuilder) copyOutput(pod *kubernetes.PodInfo, output string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pinger := watchdog.NewPinger(cb.kube, pod.Name)
	go func() {
		if err := pinger.Run(ctx); err != nil {
			fmt.Println(err)
		}
	}()

	fmt.Println("Retrieve output from diagnostics job")
	if err := cb.kube.CopyFromPod(pod.Name, output, cb.CommonOptions.Output); err != nil {
		return errors.Wrap(err, "Error while retrieving diagnostics job output")
	}
	fmt.Printf("Result successfully written to %s\n", cb.CommonOptions.Output)
	return nil
}

func (cb *CommandBuilder) storeOutputOnHost(pod *kubernetes.PodInfo, output string) error {
	hostOutput := fmt.Sprintf("%s/%s", cb.CommonOptions.OutputHostPath, output)
	fmt.Printf("Output located on host: %s, at path: %s\n", pod.Node, hostOutput)
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
		jobSpec.WithHostTmpVolume(cb.CommonOptions.OutputHostPath)
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
	outputHandler := cb.copyOutput
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
