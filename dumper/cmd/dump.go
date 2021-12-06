package cmd

import (
	"github.com/dodopizza/kubectl-shovel/internal/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	dotnetDumpBinary = "dotnet-dump"
)

type dumpOptions struct {
	*flags.DumpFlagSet
}

func newDumpCommand() *cobra.Command {
	options := &dumpOptions{}

	cmd := &cobra.Command{
		Use:   "dump [flags]",
		Args:  cobra.NoArgs,
		Short: "",
		Long:  "",
		RunE: func(cmd *cobra.Command, args []string) error {
			return makeDump(options)
		},
	}

	cmd.
		PersistentFlags().
		AddFlagSet(
			options.parse(),
		)

	return cmd
}

func (options *dumpOptions) parse() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet("dump", pflag.ExitOnError)

	options.DumpFlagSet = flags.NewDumpFlagSet()
	flagSet.AddFlagSet(options.DumpFlagSet.Parse())

	return flagSet
}

func makeDump(options *dumpOptions) error {
	args := append(
		[]string{"collect"},
		options.Args()...,
	)

	return launch(
		dotnetDumpBinary,
		args...,
	)
}
