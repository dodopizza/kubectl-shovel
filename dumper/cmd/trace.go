package cmd

import (
	"github.com/dodopizza/kubectl-shovel/internal/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	dotnetTraceBinary = "dotnet-trace"
)

type traceOptions struct {
	*flags.TraceFlagSet
}

func newTraceCommand() *cobra.Command {
	options := &traceOptions{}
	cmd := &cobra.Command{
		Use:   "trace [flags]",
		Args:  cobra.NoArgs,
		Short: "",
		Long:  "",
		RunE: func(cmd *cobra.Command, args []string) error {
			return makeTrace(options)
		},
	}

	cmd.
		PersistentFlags().
		AddFlagSet(
			options.parseFlags(),
		)

	return cmd
}

func (options *traceOptions) parseFlags() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet("trace", pflag.ExitOnError)

	options.TraceFlagSet = flags.NewTraceFlagSet()
	flagSet.AddFlagSet(options.TraceFlagSet.Parse())

	return flagSet
}

func makeTrace(options *traceOptions) error {
	args := append(
		[]string{"collect"},
		options.Args()...,
	)
	return launch(
		dotnetTraceBinary,
		args...,
	)
}
