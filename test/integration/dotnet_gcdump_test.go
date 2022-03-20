//go:build integration
// +build integration

package integration_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/dodopizza/kubectl-shovel/cli/cmd"
)

func Test_GCDumpSubcommand(t *testing.T) {
	command := "gcdump"
	testCases := cases(
		TestCase{
			name: "Custom timeout",
			args: map[string]string{
				"timeout": "60",
			},
			pod: singleContainerPod(),
		},
		TestCase{
			name: "Custom timeout with unit",
			args: map[string]string{
				"timeout": "1m",
			},
			pod: singleContainerPod(),
		},
	)

	t.Cleanup(testSetup(t, command))
	t.Run(command, func(t *testing.T) {
		for _, tc := range testCases {
			tc := tc

			t.Run(tc.name, func(t *testing.T) {
				t.Cleanup(testCaseSetup(t, &tc, command))

				args := tc.FormatArgs(command)
				shovel := cmd.NewShovelCommand()
				shovel.SetArgs(args)

				t.Logf("Execute shovel command with args: %s", args)
				err := shovel.Execute()

				require.NoError(t, err)
				if !tc.hostOutput {
					file, err := os.Stat(tc.output)
					require.NoError(t, err)
					require.NotEmpty(t, file.Size())
				}
			})
		}
	})
}
