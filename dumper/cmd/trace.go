package cmd

import (
	"github.com/dodopizza/kubectl-shovel/internal/flags"
	"github.com/spf13/cobra"
)

func newTraceCommand() *cobra.Command {
	options := NewDiagnosticBinaryOptions(
		"dotnet-trace",
		"trace",
		flags.NewTraceFlagSet,
	)
	return options.GetCommand()
}
