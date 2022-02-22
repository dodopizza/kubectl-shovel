//go:build integration
// +build integration

package integration_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/dodopizza/kubectl-shovel/cli/cmd"
)

func Test_DumpSubcommand(t *testing.T) {
	testCases := cases(
		TestCase{
			name: "Custom type",
			args: []string{
				"--type", "Heap",
			},
			pod: singleContainerPod(),
		},
	)

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			teardown := setup(t, tc, "dump-test")
			defer teardown()

			c := cmd.NewShovelCommand()
			c.SetArgs(tc.FormatArgs("dump"))
			require.NoError(t, c.Execute())

			if !tc.hostOutput {
				file, err := os.Stat(tc.output)
				require.NoError(t, err)
				require.NotEmpty(t, file.Size())
			}
		})
	}
}
