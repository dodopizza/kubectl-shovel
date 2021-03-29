package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/dodopizza/kubectl-shovel/internal/version"
)

type commonOptions struct {
	image   string
	podName string
	output  string

	kubeFlags *genericclioptions.ConfigFlags
}

func (options *commonOptions) newCommonFlags(tool string) *pflag.FlagSet {
	flags := pflag.NewFlagSet("common", pflag.ExitOnError)
	flags.StringVar(
		&options.podName,
		"pod-name",
		options.podName,
		"Target pod",
	)
	panicOnError(cobra.MarkFlagRequired(flags, "pod-name"))

	flags.StringVarP(
		&options.output,
		"output",
		"o",
		"./"+
			currentTime()+
			"."+
			tool,
		"Output file",
	)

	flags.StringVar(
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
	options.kubeFlags.AddFlags(flags)

	return flags
}
