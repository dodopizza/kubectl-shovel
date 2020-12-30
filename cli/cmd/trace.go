package cmd

import (
	"fmt"
	"strings"

	"github.com/dodopizza/kubectl-shovel/internal/version"
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
		Long: "This subcommand will capture 10 seconds of runtime events with dotnet-trace tool for running in k8s appplication.\n" +
			"Result will be saved locally in nettrace format so you'll be able to convert it and analyze with appropriate tools.\n" +
			"You can find more info about dotnet-trace tool by the following links:\n\n" +
			"\t* https://github.com/dotnet/diagnostics/blob/master/documentation/dotnet-trace-instructions.md\n" +
			"\t* https://docs.microsoft.com/en-us/dotnet/core/diagnostics/dotnet-trace",
		Example: fmt.Sprintf(examplesTemplate, "trace") + "\n\n" +
			"One of the resulting trace usage examples is converting it to speedscope format:\n\n" +
			"\tdotnet trace convert myapp.trace --format Speedscope\n\n" +
			"And then analyzing it with https://www.speedscope.app/",
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
		"trace",
	)
}
