package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dodopizza/kubectl-shovel/cli/kubernetes"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type gcDumpOptions struct {
	podName string
	output  string

	kubeFlags *genericclioptions.ConfigFlags
}

func newGCDumpCommand() *cobra.Command {
	options := &gcDumpOptions{}
	cmd := &cobra.Command{
		Use:   "gcdump [flags]",
		Short: "Get dotnet-gcdump results",
		RunE: func(*cobra.Command, []string) error {
			return options.makeGCDump()
		},
	}

	cmd.
		PersistentFlags().
		AddFlagSet(
			options.checkFlags(),
		)

	return cmd
}

func (options *gcDumpOptions) checkFlags() *pflag.FlagSet {
	flags := pflag.NewFlagSet("gcdump", pflag.ExitOnError)
	flags.StringVar(&options.podName, "pod-name", options.podName, "Pod name for creating dump")
	panicOnError(cobra.MarkFlagRequired(flags, "pod-name"))

	flags.StringVar(
		&options.output,
		"output",
		"./"+
			strconv.Itoa(
				int(time.Now().Unix()),
			)+
			".gcdump",
		"Dump output file",
	)

	options.kubeFlags = genericclioptions.NewConfigFlags(false)
	options.kubeFlags.AddFlags(flags)

	return flags
}

func (options *gcDumpOptions) makeGCDump() error {
	k8s, err := kubernetes.NewClient(options.kubeFlags)
	if err != nil {
		return nil
	}
	pod, err := k8s.GetPodInfo(options.podName)
	if err != nil {
		return err
	}

	jobName, err := newJobName()
	if err != nil {
		return err
	}
	containerID := strings.TrimPrefix(pod.Status.ContainerStatuses[0].ContainerID, "docker://")
	fmt.Printf(
		"Run job %s for pod %s with container ID %s\n",
		jobName,
		options.podName,
		containerID,
	)
	err = k8s.RunJob(
		jobName,
		dumperImageName,
		pod.Spec.NodeName,
		[]string{
			"--container-id",
			containerID,
		},
	)
	if err != nil {
		return err
	}

	return nil
}
