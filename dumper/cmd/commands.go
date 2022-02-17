package cmd

import (
	"github.com/dodopizza/kubectl-shovel/internal/flags"
	"github.com/spf13/cobra"
)

func newTraceCommand() *cobra.Command {
	builder := NewCommandBuilder(
		"dotnet-trace",
		"trace",
		flags.NewDotnetTrace,
	)
	return builder.Build()
}

func newGCDumpCommand() *cobra.Command {
	builder := NewCommandBuilder(
		"dotnet-gcdump",
		"gcdump",
		flags.NewDotnetGCDump,
	)
	return builder.Build()
}

func newDumpCommand() *cobra.Command {
	options := NewCommandBuilder(
		"dotnet-dump",
		"dump",
		flags.NewDotnetDump,
	)
	return options.Build()
}
