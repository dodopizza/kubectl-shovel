package cmd

import (
	"github.com/dodopizza/kubectl-shovel/internal/flags"
	"github.com/spf13/cobra"
)

func newDumpCommand() *cobra.Command {
	options := NewDiagnosticBinaryOptions(
		"dotnet-dump",
		"dump",
		flags.NewDumpFlagSet,
	)
	return options.GetCommand()
}
