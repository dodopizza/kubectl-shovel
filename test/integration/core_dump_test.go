//go:build integration
// +build integration

package integration_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/dodopizza/kubectl-shovel/cli/cmd"
)

func Test_CoreDumpSubcommand(t *testing.T) {
	testCases := cases(
		TestCase{
			name: "Custom type",
			args: map[string]string{
				"type": "Heap",
			},
			pod: singleContainerPod(),
		},
	)

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			teardown := setup(t, &tc, "coredump-test")
			defer teardown()
			args := tc.FormatArgs("coredump")
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
}
