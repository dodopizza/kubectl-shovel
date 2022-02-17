package cmd

import (
	"fmt"
	"github.com/dodopizza/kubectl-shovel/internal/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// todo: better naming

type DiagnosticBinaryOptions struct {
	BinaryName              string
	FlagSetContainer        flags.FlagSetContainer
	FlagSetContainerFactory flags.FlagSetContainerFactory
	ToolName                string
}

func NewDiagnosticBinaryOptions(binary, tool string, factory flags.FlagSetContainerFactory) *DiagnosticBinaryOptions {
	return &DiagnosticBinaryOptions{
		BinaryName:              binary,
		FlagSetContainerFactory: factory,
		ToolName:                tool,
	}
}

func (db *DiagnosticBinaryOptions) GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("%s [flags]", db.ToolName),
		Args:  cobra.NoArgs,
		Short: "",
		Long:  "",
		RunE: func(cmd *cobra.Command, args []string) error {
			return db.Run()
		},
	}

	cmd.PersistentFlags().AddFlagSet(db.GetFlags())
	return cmd
}

func (db *DiagnosticBinaryOptions) GetFlags() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet(db.ToolName, pflag.ExitOnError)

	db.FlagSetContainer = db.FlagSetContainerFactory()
	flagSet.AddFlagSet(db.FlagSetContainer.Parse())

	return flagSet
}

func (db *DiagnosticBinaryOptions) Run() error {
	args := append(
		[]string{"collect"},
		db.FlagSetContainer.Args()...,
	)

	return launch(
		db.BinaryName,
		args...,
	)
}
