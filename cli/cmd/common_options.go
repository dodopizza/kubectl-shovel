package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/dodopizza/kubectl-shovel/internal/version"
)

type commonOptions struct {
	image         string
	podName       string
	output        string
	containerName string

	kubeFlags *genericclioptions.ConfigFlags
}

func (options *commonOptions) newCommonFlags(tool string) *pflag.FlagSet {
	flagSet := pflag.NewFlagSet("common", pflag.ExitOnError)
	flagSet.StringVar(
		&options.podName,
		"pod-name",
		options.podName,
		"Target pod",
	)
	panicOnError(cobra.MarkFlagRequired(flagSet, "pod-name"))

	flagSet.StringVarP(
		&options.containerName,
		"container",
		"c",
		options.containerName,
		"Target container in pod. Required if pod run multiple containers",
	)

	flagSet.StringVarP(
		&options.output,
		"output",
		"o",
		"./"+
			currentTime()+
			"."+
			tool,
		"Output file",
	)

	flagSet.StringVar(
		&options.image,
		"image",
		strings.Join(
			[]string{
				dumperImageName,
				version.GetVersion(),
			},
			":",
		),
		"Image of dumper to use for job",
	)
	options.kubeFlags = genericclioptions.NewConfigFlags(false)
	options.kubeFlags.AddFlags(flagSet)

	return flagSet
}
