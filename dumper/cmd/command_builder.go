package cmd

import (
	"fmt"
	"github.com/dodopizza/kubectl-shovel/internal/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// CommandBuilder represents logic for building and running tools
type CommandBuilder struct {
	BinaryName              string
	FlagSetContainer        flags.FlagSetContainer
	FlagSetContainerFactory flags.FlagSetContainerFactory
	ToolName                string
}

// NewCommandBuilder returns options with specified tool name
func NewCommandBuilder(binary, tool string, factory flags.FlagSetContainerFactory) *CommandBuilder {
	return &CommandBuilder{
		BinaryName:              binary,
		FlagSetContainerFactory: factory,
		ToolName:                tool,
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
	flagSet := pflag.NewFlagSet(db.ToolName, pflag.ExitOnError)

	db.FlagSetContainer = db.FlagSetContainerFactory()
	flagSet.AddFlagSet(db.FlagSetContainer.GetFlags())

	return flagSet
}

func (db *CommandBuilder) run() error {
	args := append(
		[]string{"collect"},
		db.FlagSetContainer.GetArgs()...,
	)

	return launch(
		db.BinaryName,
		args...,
	)
}
