//go:build integration
// +build integration

package integration_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	core "k8s.io/api/core/v1"

	"github.com/dodopizza/kubectl-shovel/cli/cmd"
)

func Test_TraceSubcommand(t *testing.T) {
	testCases := []struct {
		name string
		args []string
		pod  *core.Pod
	}{
		{
			name: "Basic test",
			args: []string{},
			pod:  singleContainerPod(),
		},
		{
			name: "Custom duration",
			args: []string{
				"--duration",
				"00:00:00:30",
			},
			pod: singleContainerPod(),
		},
		{
			name: "Custom duration with units",
			args: []string{
				"--duration",
				"1m",
			},
			pod: singleContainerPod(),
		},
		{
			name: "Custom format",
			args: []string{
				"--format",
				"Speedscope",
			},
			pod: singleContainerPod(),
		},
		{
			name: "Multicontainer pod",
			args: []string{
				"--container",
				targetContainerName,
			},
			pod: multiContainerPod(),
		},
		{
			name: "MultiContainer pod with shared mount",
			args: []string{
				"--container",
				targetContainerName,
			},
			pod: multiContainerPodWithSharedMount(),
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

			c := cmd.NewShovelCommand()
			c.SetArgs(args)
			require.NoError(t, c.Execute())

			file, err := os.Stat(outputFilename)
			require.NoError(t, err)
			require.NotEmpty(t, file.Size())
		})
	}
}
