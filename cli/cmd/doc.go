package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func newDocCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "doc",
		Hidden: true,
		Args:   cobra.NoArgs,
		Short:  "Generate documentation",
		Long:   "This command will generate documentation for this CLI.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doc.GenMarkdownTree(cmd.Root(), "./docs")
		},
	}

	return cmd
}
