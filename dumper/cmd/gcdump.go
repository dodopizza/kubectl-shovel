package cmd

import (
	"github.com/dodopizza/kubectl-shovel/internal/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	dotnetGCDumpBinary = "dotnet-gcdump"
)

type gcDumpOptions struct {
	*flags.GCDumpFlagSet
}

func newGCDumpCommand() *cobra.Command {
	options := &gcDumpOptions{}
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
	flagSet := pflag.NewFlagSet("gcdump", pflag.ExitOnError)

	options.GCDumpFlagSet = flags.NewGCDumpFlagSet()
	flagSet.AddFlagSet(options.GCDumpFlagSet.Parse())

	return flagSet
}

func makeGCDump(options *gcDumpOptions) error {
	args := append(
		[]string{"collect"},
		options.Args()...,
	)
	return launch(
		dotnetGCDumpBinary,
		args...,
	)
}
