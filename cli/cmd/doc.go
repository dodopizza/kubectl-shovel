package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// NewDocCommand return command that generate tool docs from sources
func NewDocCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "doc",
		Hidden: true,
		Args:   cobra.NoArgs,
		Short:  "Generate documentation",
		Long:   "This command will generate documentation for this CLI.",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return doc.GenMarkdownTree(cmd.Root(), "./docs")
		},
	}

	return cmd
}
