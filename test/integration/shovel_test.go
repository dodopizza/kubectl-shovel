// +build integration

package integration_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/dodopizza/kubectl-shovel/cli/cmd"
)

var (
	namespace      = "default"
	cliPath        = "./bin/kubectl-shovel"
	dumperImage    = "kubectl-shovel/dumper-integration-tests"
	tempDirPattern = "*-kubectl-shovel"
)

func Test_TraceSubcommand(t *testing.T) {
	sampleAppName, teardown := setup(t)
	defer teardown()
	dir, _ := ioutil.TempDir("", tempDirPattern)
	defer os.RemoveAll(dir)
	outputFilename := filepath.Join(dir, "trace-test")

	args := []string{
		"trace",
		"--pod-name",
		sampleAppName,
		"--output",
		outputFilename,
		"--image",
		dumperImage,
		"--container",
		"target",
	}
	cmd := cmd.NewShovelCommand()
	cmd.SetArgs(args)
	require.NoError(t, cmd.Execute())

	file, err := os.Stat(outputFilename)
	require.NoError(t, err)
	require.NotEmpty(t, file.Size())
}

func Test_GCDumpSubcommand(t *testing.T) {
	targetPodName, teardown := setup(t)
	defer teardown()
	dir, _ := ioutil.TempDir("", tempDirPattern)
	defer os.RemoveAll(dir)
	outputFilename := filepath.Join(dir, "gcdump-test")
	args := []string{
		"gcdump",
		"--pod-name",
		targetPodName,
		"--output",
		outputFilename,
		"--image",
		dumperImage,
		"--container",
		"target",
	}
	cmd := cmd.NewShovelCommand()
	cmd.SetArgs(args)
	require.NoError(t, cmd.Execute())

	file, err := os.Stat(outputFilename)
	require.NoError(t, err)
	require.NotEmpty(t, file.Size())
}
