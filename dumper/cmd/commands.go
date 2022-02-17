package cmd

import (
	"github.com/dodopizza/kubectl-shovel/internal/flags"
	"github.com/spf13/cobra"
)

func newTraceCommand() *cobra.Command {
	builder := NewCommandBuilder(
		"dotnet-trace",
		"trace",
		flags.NewTraceFlagSet,
	)
	return builder.Build()
}

func newGCDumpCommand() *cobra.Command {
	builder := NewCommandBuilder(
		"dotnet-gcdump",
		"gcdump",
		flags.NewGCDumpFlagSet,
	)
	return builder.Build()
}

func newDumpCommand() *cobra.Command {
	options := NewCommandBuilder(
		"dotnet-dump",
		"dump",
		flags.NewDumpFlagSet,
	)
	return options.Build()
}
