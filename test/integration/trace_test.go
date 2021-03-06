// +build integration

package integration_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"

	"github.com/dodopizza/kubectl-shovel/cli/cmd"
)

func Test_TraceSubcommand(t *testing.T) {
	testCases := []struct {
		name string
		args []string
		pod  *v1.Pod
	}{
		{
			name: "Basic test",
			args: []string{},
			pod:  sampleAppPod(),
		},
		{
			name: "Custom duration",
			args: []string{
				"--duration",
				"00:00:00:30",
			},
			pod: sampleAppPod(),
		},
		{
			name: "Custom duration with units",
			args: []string{
				"--duration",
				"1m",
			},
			pod: sampleAppPod(),
		},
		{
			name: "Custom format",
			args: []string{
				"--format",
				"Speedscope",
			},
			pod: sampleAppPod(),
		},
		{
			name: "Multicontainer pod",
			args: []string{
				"--container",
				targetContainerName,
			},
			pod: multiContainerPod(),
		},
	}
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
			cmd := cmd.NewShovelCommand()
			cmd.SetArgs(args)
			require.NoError(t, cmd.Execute())

			file, err := os.Stat(outputFilename)
			require.NoError(t, err)
			require.NotEmpty(t, file.Size())
		})
	}
}
