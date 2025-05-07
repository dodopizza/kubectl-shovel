//go:build integration
// +build integration

package integration_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	
	"github.com/dodopizza/kubectl-shovel/cli/cmd"
)

func Test_InitContainer_Support(t *testing.T) {
	defer testSetup(t, "init_container_test")()

	for _, tc := range []*TestCase{
		NewTestCase("Init container explicitly specified").
			WithPod(podWithInitContainer()).
			WithArgs("container", initContainerName).
			DownloadOutput(),
		NewTestCase("Init container as side container").
			WithPod(podWithInitContainerSidecar()).
			WithArgs("container", "side-container").
			DownloadOutput(),
	} {
		t.Run(tc.name, func(t *testing.T) {
			defer testCaseSetup(t, tc, "init_container_test")()

			args := tc.FormatArgs("dump")

			t.Logf("Run command with args: %v", args)
			output, err := execCommand(t, args)
			
			t.Logf("Command output: %s", output)
			require.NoError(t, err)
			
			if !tc.hostOutput {
				_, err = os.Stat(tc.output)
				require.NoError(t, err)
			}
		})
	}
}