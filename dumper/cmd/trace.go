package cmd

import (
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	dotnetTraceBinary = "dotnet-trace"
)

type traceOptions struct {
	pid    int
	output string
}

func newTraceOptions() *traceOptions {
	return &traceOptions{
		pid: 1,
	}
}

func newTraceCommand() *cobra.Command {
	options := newTraceOptions()
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
	flags := pflag.NewFlagSet("trace", pflag.ExitOnError)

	return flags
}

func makeTrace(options *traceOptions) error {
	return launch(
		dotnetTraceBinary,
		"collect",
		"--process-id",
		strconv.Itoa(options.pid),
		"--duration",
		"00:00:00:10",
	)
}
