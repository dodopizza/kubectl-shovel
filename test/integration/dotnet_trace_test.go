//go:build integration
// +build integration

package integration_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/dodopizza/kubectl-shovel/cli/cmd"
)

func Test_TraceSubcommand(t *testing.T) {
	command := "trace"
	testCases := cases(
		NewTestCase("Custom duration").WithArgs("duration", "00:00:00:30"),
		NewTestCase("Custom duration with units").WithArgs("duration", "1m"),
		NewTestCase("Custom format").WithArgs("format", "Speedscope"),
	)

	t.Cleanup(testSetup(t, command))
	t.Run(command, func(t *testing.T) {
		for _, tc := range testCases {
			tc := tc

			t.Run(tc.name, func(t *testing.T) {
				t.Cleanup(testCaseSetup(t, tc, command))

				args := tc.FormatArgs(command)
				shovel := cmd.NewShovelCommand()
				shovel.SetArgs(args)

				t.Logf("Execute shovel command with args: %s", args)
				err := shovel.Execute()

				require.NoError(t, err)
				if !tc.hostOutput {
					t.Logf("Looking for artifact at path: %s", tc.output)
					file, err := os.Stat(tc.output)
					require.NoError(t, err)
					require.NotEmpty(t, file.Size())
				}
			})
		}
	})
}
