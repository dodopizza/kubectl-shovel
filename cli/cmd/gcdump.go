package cmd

import (
	"strings"

	"github.com/dodopizza/kubectl-shovel/pkg/version"
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
		"dotnet-gcdump",
	)
}
