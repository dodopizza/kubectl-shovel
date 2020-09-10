package cmd

import (
	"strings"

	"github.com/dodopizza/kubectl-shovel/pkg/version"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type traceOptions struct {
	image   string
	podName string
	output  string

	kubeFlags *genericclioptions.ConfigFlags
}

func newTraceCommand() *cobra.Command {
	options := &traceOptions{}
	cmd := &cobra.Command{
		Use:   "trace [flags]",
		Short: "Get dotnet-trace results",
		RunE: func(*cobra.Command, []string) error {
			return options.maketrace()
		},
	}

	cmd.
		PersistentFlags().
		AddFlagSet(
			options.checkFlags(),
		)

	return cmd
}

func (options *traceOptions) checkFlags() *pflag.FlagSet {
	flags := pflag.NewFlagSet("trace", pflag.ExitOnError)
	flags.StringVar(&options.podName, "pod-name", options.podName, "Pod name for creating dump")
	panicOnError(cobra.MarkFlagRequired(flags, "pod-name"))

	flags.StringVarP(
		&options.output,
		"output",
		"o",
		"./"+
			currentTime()+
			".trace",
		"Dump output file",
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

func (options *traceOptions) maketrace() error {
	return run(
		options.kubeFlags,
		options.image,
		options.podName,
		options.output,
		"dotnet-trace",
	)
}
