package cmd

import (
	"github.com/dodopizza/kubectl-shovel/internal/flags"
	"github.com/spf13/cobra"
)

// NewGCDumpCommand return command that perform dotnet-gcdump on target process
func NewGCDumpCommand(commonOptions *ContainerOptions) *cobra.Command {
	return NewCommandBuilder(commonOptions, flags.NewDotnetGCDump).Build()
}

// NewTraceCommand return command that perform dotnet-trace on target process
func NewTraceCommand(commonOptions *ContainerOptions) *cobra.Command {
	return NewCommandBuilder(commonOptions, flags.NewDotnetTrace).Build()
}

// NewDumpCommand return command that perform dotnet-dump on target process
func NewDumpCommand(commonOptions *ContainerOptions) *cobra.Command {
	return NewCommandBuilder(commonOptions, flags.NewDotnetDump).Build()
}
