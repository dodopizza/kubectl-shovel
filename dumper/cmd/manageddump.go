package cmd

import (
	"github.com/dodopizza/kubectl-shovel/internal/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	dotnetDumpBinary = "dotnet-dump"
)

type managedDumpOptions struct {
	*flags.ManagedDumpFlagSet
}

func newManagedDumpCommand() *cobra.Command {
	options := &managedDumpOptions{}

	cmd := &cobra.Command{
		Use:   "dump [flags]",
		Args:  cobra.NoArgs,
		Short: "",
		Long:  "",
		RunE: func(cmd *cobra.Command, args []string) error {
			return makeManagedDump(options)
		},
	}

	cmd.
		PersistentFlags().
		AddFlagSet(
			options.parse(),
		)

	return cmd
}

func (options *managedDumpOptions) parse() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet("dump", pflag.ExitOnError)

	options.ManagedDumpFlagSet = flags.NewManagedDumpFlagSet()
	flagSet.AddFlagSet(options.ManagedDumpFlagSet.Parse())

	return flagSet
}

func makeManagedDump(options *managedDumpOptions) error {
	args := append(
		[]string{"collect"},
		options.Args()...,
	)

	return launch(
		dotnetDumpBinary,
		args...,
	)
}
