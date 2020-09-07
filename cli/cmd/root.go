package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var podName string
var output string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dotnet-k8s-dumper",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		podInfo, err := getPodInfo(podName)

		if err != nil {
			return err
		}

		log.Println(podInfo)

		return nil
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&podName, "pod-name", "", "Pod name for creating dump")
	rootCmd.PersistentFlags().StringVar(&output, "output", "", "Dump output file")
}
