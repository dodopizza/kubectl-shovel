package cmd

import (
	"github.com/spf13/cobra"

	"github.com/dodopizza/kubectl-shovel/internal/flags"
)

// NewGCDumpCommand return command that perform dotnet-gcdump on target process
func NewGCDumpCommand(commonOptions *CommonOptions) *cobra.Command {
	return NewCommandBuilder(commonOptions, flags.NewDotnetGCDump).Build()
}

// NewTraceCommand return command that perform dotnet-trace on target process
func NewTraceCommand(commonOptions *CommonOptions) *cobra.Command {
	return NewCommandBuilder(commonOptions, flags.NewDotnetTrace).Build()
}

// NewDumpCommand return command that perform dotnet-dump on target process
func NewDumpCommand(commonOptions *CommonOptions) *cobra.Command {
	return NewCommandBuilder(commonOptions, flags.NewDotnetDump).Build()
}

// NewCoreDumpCommand return command that perform createdump on target process
func NewCoreDumpCommand(commonOptions *CommonOptions) *cobra.Command {
	return NewCommandBuilder(commonOptions, flags.NewCoreDump).Build()
}
