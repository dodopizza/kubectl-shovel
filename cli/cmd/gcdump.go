package cmd

import (
	"fmt"
	"strings"

	"github.com/dodopizza/kubectl-shovel/internal/version"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type gcDumpOptions struct {
	image   string
	podName string
	output  string

	kubeFlags *genericclioptions.ConfigFlags
}

func newGCDumpCommand() *cobra.Command {
	options := &gcDumpOptions{}
	cmd := &cobra.Command{
		Use:   "gcdump [flags]",
		Short: "Get dotnet-gcdump results",
		Long: "This subcommand will run dotnet-gcdump tool for running in k8s appplication.\n" +
			"Result will be saved locally so you'll be able to analyze it with appropriate tools.\n" +
			"You can find more info about dotnet-gcdump tool by the following links:\n\n" +
			"\t* https://devblogs.microsoft.com/dotnet/collecting-and-analyzing-memory-dumps/\n" +
			"\t* https://docs.microsoft.com/en-us/dotnet/core/diagnostics/dotnet-gcdump",
		Example: fmt.Sprintf(examplesTemplate, "gcdump"),
		RunE: func(*cobra.Command, []string) error {
			return options.makeGCDump()
		},
	}

	cmd.
		PersistentFlags().
		AddFlagSet(
			options.parseFlags(),
		)

	return cmd
}

func (options *gcDumpOptions) parseFlags() *pflag.FlagSet {
	flags := pflag.NewFlagSet("gcdump", pflag.ExitOnError)
	flags.StringVar(&options.podName, "pod-name", options.podName, "Pod name for creating dump")
	panicOnError(cobra.MarkFlagRequired(flags, "pod-name"))

	flags.StringVarP(
		&options.output,
		"output",
		"o",
		"./"+
			currentTime()+
			".gcdump",
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

func (options *gcDumpOptions) makeGCDump() error {
	return run(
		options.kubeFlags,
		options.image,
		options.podName,
		options.output,
		"gcdump",
	)
}
