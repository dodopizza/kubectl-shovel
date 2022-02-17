package cmd

import (
	"fmt"
	"github.com/dodopizza/kubectl-shovel/internal/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// CommandBuilder represents logic for building and running tools
type CommandBuilder struct {
	BinaryName string
	FlagSet    flags.DotnetTool
	ToolName   string
}

// NewCommandBuilder returns options with specified tool name
func NewCommandBuilder(binary, tool string, factory flags.DotnetToolFactory) *CommandBuilder {
	return &CommandBuilder{
		BinaryName: binary,
		FlagSet:    factory(),
		ToolName:   tool,
	}
}

// Build will build command *cobra.Command from options
func (db *CommandBuilder) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("%s [flags]", db.ToolName),
		Args:  cobra.NoArgs,
		Short: "",
		Long:  "",
		RunE: func(cmd *cobra.Command, args []string) error {
			return db.run()
		},
	}

	cmd.PersistentFlags().AddFlagSet(db.flags())
	return cmd
}

func (db *CommandBuilder) flags() *pflag.FlagSet {
	fs := pflag.NewFlagSet(db.ToolName, pflag.ExitOnError)
	fs.AddFlagSet(db.FlagSet.GetFlags())
	return fs
}

func (db *CommandBuilder) run() error {
	args := append(
		[]string{"collect"},
		db.FlagSet.GetArgs()...,
	)

	return launch(
		db.BinaryName,
		args...,
	)
}
