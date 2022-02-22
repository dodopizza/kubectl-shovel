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
	testCases := cases(
		TestCase{
			name: "Custom duration",
			args: []string{
				"--duration",
				"00:00:00:30",
			},
			pod: singleContainerPod(),
		},
		TestCase{
			name: "Custom duration with units",
			args: []string{
				"--duration",
				"1m",
			},
			pod: singleContainerPod(),
		},
		TestCase{
			name: "Custom format",
			args: []string{
				"--format",
				"Speedscope",
			},
			pod: singleContainerPod(),
		},
	)

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			teardown := setup(t, tc, "trace-test")
			defer teardown()

			c := cmd.NewShovelCommand()
			c.SetArgs(tc.FormatArgs("trace"))
			require.NoError(t, c.Execute())

			if !tc.hostOutput {
				file, err := os.Stat(tc.output)
				require.NoError(t, err)
				require.NotEmpty(t, file.Size())
			}
		})
	}
}
