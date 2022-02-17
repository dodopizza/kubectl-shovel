package cmd

import (
	"github.com/spf13/cobra"
)

func NewDumperCommand() *cobra.Command {
	options := &commonOptions{}
	cmd := &cobra.Command{
		Use:               "dumper",
		Short:             "Tool to gather diagnostic information from dotnet process",
		SilenceUsage:      true,
		DisableAutoGenTag: true,
	}
	cmd.
		PersistentFlags().
		AddFlagSet(options.GetFlags())
	cmd.AddCommand(NewGCDumpCommand(options))
	cmd.AddCommand(NewTraceCommand(options))
	cmd.AddCommand(NewDumpCommand(options))

	return cmd
}
