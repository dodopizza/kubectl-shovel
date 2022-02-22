//go:build integration
// +build integration

package integration_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/dodopizza/kubectl-shovel/cli/cmd"
	"github.com/stretchr/testify/require"
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
			teardown := setup(t, tc.pod)
			defer teardown()
			dir, _ := ioutil.TempDir("", tempDirPattern)
			defer os.RemoveAll(dir)
			outputFilename := filepath.Join(dir, "dump-test")
			args := append([]string{
				"dump",
				"--pod-name",
				tc.pod.Name,
				"--output",
				outputFilename,
				"--image",
				dumperImage,
			}, tc.args...)

			c := cmd.NewShovelCommand()
			c.SetArgs(args)
			require.NoError(t, c.Execute())

			if !tc.hostOutput {
				file, err := os.Stat(outputFilename)
				require.NoError(t, err)
				require.NotEmpty(t, file.Size())
			}
		})
	}
}
