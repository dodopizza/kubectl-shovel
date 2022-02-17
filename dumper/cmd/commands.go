package cmd

import (
	"github.com/dodopizza/kubectl-shovel/internal/flags"
	"github.com/spf13/cobra"
)

// NewGCDumpCommand return command that perform dotnet-gcdump on target process
func NewGCDumpCommand() *cobra.Command {
	return NewCommandBuilder(flags.NewDotnetGCDump).Build()
}

// NewTraceCommand return command that perform dotnet-trace on target process
func NewTraceCommand() *cobra.Command {
	return NewCommandBuilder(flags.NewDotnetTrace).Build()
}

// NewDumpCommand return command that perform dotnet-dump on target process
func NewDumpCommand() *cobra.Command {
	return NewCommandBuilder(flags.NewDotnetDump).Build()
}
