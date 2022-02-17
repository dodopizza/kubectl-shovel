package cmd

import (
	"github.com/dodopizza/kubectl-shovel/internal/flags"
	"github.com/spf13/cobra"
)

func newGCDumpCommand() *cobra.Command {
	options := NewDiagnosticBinaryOptions(
		"dotnet-gcdump",
		"gcdump",
		flags.NewGCDumpFlagSet,
	)
	return options.GetCommand()
}
