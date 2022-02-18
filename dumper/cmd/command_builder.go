package cmd

import (
	"fmt"
	"github.com/dodopizza/kubectl-shovel/internal/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type commonOptions struct {
	containerID      string
	containerRuntime string
}

// CommandBuilder represents logic for building and running tools
type CommandBuilder struct {
	CommonOptions *commonOptions
	tool          flags.DotnetTool
}

func (options *commonOptions) GetFlags() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet("common", pflag.ExitOnError)
	flagSet.StringVar(
		&options.containerID,
		"container-id",
		options.containerID,
		"Container ID to run tool for",
	)
	flagSet.StringVar(
		&options.containerRuntime,
		"container-runtime",
		options.containerRuntime,
		"Container Runtime to run tool for",
	)
	_ = cobra.MarkFlagRequired(flagSet, "container-id")
	_ = cobra.MarkFlagRequired(flagSet, "container-runtime")

	return flagSet
}

// NewCommandBuilder returns options with specified tool name
func NewCommandBuilder(options *commonOptions, factory flags.DotnetToolFactory) *CommandBuilder {
	return &CommandBuilder{
		tool:          factory(),
		CommonOptions: options,
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
