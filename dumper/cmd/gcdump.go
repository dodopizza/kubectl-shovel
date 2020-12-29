package cmd

import (
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	dotnetGCDumpBinary = "dotnet-gcdump"
)

type gcDumpOptions struct {
	pid int
}

func newGCDumpOptions() *gcDumpOptions {
	return &gcDumpOptions{
		pid: 1,
	}
}

func newGCDumpCommand() *cobra.Command {
	options := newGCDumpOptions()
	cmd := &cobra.Command{
		Use:   "gcdump [flags]",
		Args:  cobra.NoArgs,
		Short: "",
		Long:  "",
		RunE: func(cmd *cobra.Command, args []string) error {
			return makeGCDump(options)
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

	return flags
}

func makeGCDump(options *gcDumpOptions) error {
	return launch(
		dotnetGCDumpBinary,
		"collect",
		"--process-id",
		strconv.Itoa(options.pid),
	)
}
