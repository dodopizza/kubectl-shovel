package cmd

import (
	"fmt"
	"github.com/dodopizza/kubectl-shovel/internal/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// ContainerOptions represents container info
type ContainerOptions struct {
	ID      string
	Runtime string
}

// CommandBuilder represents logic for building and running tools
type CommandBuilder struct {
	ContainerOptions *ContainerOptions
	tool             flags.DotnetTool
}

// GetFlags return FlagSet that describes options for container selection
func (options *ContainerOptions) GetFlags() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet("common", pflag.ExitOnError)
	flagSet.StringVar(
		&options.ID,
		"container-id",
		options.ID,
		"Container ID to run tool for",
	)
	flagSet.StringVar(
		&options.Runtime,
		"container-runtime",
		options.Runtime,
		"Container Runtime to run tool for",
	)
	_ = cobra.MarkFlagRequired(flagSet, "container-id")
	_ = cobra.MarkFlagRequired(flagSet, "container-runtime")

	return flagSet
}

// NewCommandBuilder returns options with specified tool name
func NewCommandBuilder(options *ContainerOptions, factory flags.DotnetToolFactory) *CommandBuilder {
	return &CommandBuilder{
		tool:             factory(),
		ContainerOptions: options,
	}
}

// Build will build command *cobra.Command from options
func (cb *CommandBuilder) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("%s [flags]", cb.tool.ToolName()),
		Args:  cobra.NoArgs,
		Short: "",
		Long:  "",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cb.launch()
		},
	}

	cmd.PersistentFlags().AddFlagSet(cb.flags())
	return cmd
}

func (cb *CommandBuilder) flags() *pflag.FlagSet {
	fs := pflag.NewFlagSet(cb.tool.ToolName(), pflag.ExitOnError)
	fs.AddFlagSet(cb.tool.GetFlags())
	return fs
}
