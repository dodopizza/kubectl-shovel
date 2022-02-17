package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/dodopizza/kubectl-shovel/internal/globals"
)

func newVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Args:  cobra.NoArgs,
		Short: "Print your cli version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(globals.GetVersion())
		},
	}

	return cmd
}
