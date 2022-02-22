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
			teardown := setup(t, tc.pod)
			defer teardown()
			dir, _ := ioutil.TempDir("", tempDirPattern)
			defer os.RemoveAll(dir)
			outputFilename := filepath.Join(dir, "trace-test")
			args := append([]string{
				"trace",
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
