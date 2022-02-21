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

func Test_GCDumpSubcommand(t *testing.T) {
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
			name: "Custom timeout",
			args: []string{
				"--timeout",
				"60",
			},
			pod: singleContainerPod(),
		},
		{
			name: "Custom timeout with unit",
			args: []string{
				"--timeout",
				"1m",
			},
			pod: singleContainerPod(),
		},
		{
			name: "MultiContainer pod",
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
			outputFilename := filepath.Join(dir, "gcdump-test")
			args := append([]string{
				"gcdump",
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
