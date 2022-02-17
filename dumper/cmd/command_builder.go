package cmd

import (
	"fmt"
	"github.com/dodopizza/kubectl-shovel/internal/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// CommandBuilder represents logic for building and running tools
type CommandBuilder struct {
	tool flags.DotnetTool
}

// NewCommandBuilder returns options with specified tool name
func NewCommandBuilder(factory flags.DotnetToolFactory) *CommandBuilder {
	return &CommandBuilder{
		tool: factory(),
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
			return cb.run()
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

func (cb *CommandBuilder) run() error {
	args := append(
		[]string{"collect"},
		cb.tool.GetArgs()...,
	)

	return launch(
		cb.tool.BinaryName(),
		args...,
	)
}
