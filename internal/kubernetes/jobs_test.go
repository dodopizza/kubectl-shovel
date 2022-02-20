package kubernetes

import (
	"fmt"
	"github.com/dodopizza/kubectl-shovel/internal/globals"
	"github.com/stretchr/testify/require"
	core "k8s.io/api/core/v1"
	"testing"
)

func Test_NewRunJobSpec(t *testing.T) {
	generator = func() string {
		return "suffix"
	}
	expJobName := fmt.Sprintf("%s-suffix", globals.PluginName)
	testCases := []struct {
		args      []string
		image     string
		name      string
		pod       *PodInfo
		container *ContainerInfo
	}{
		{
			name:  "Correct job spec generates",
			args:  []string{"/bin/sh", "-c", "sleep", "1"},
			image: "alpine",
			pod: NewPodInfo(&core.Pod{
				Spec: core.PodSpec{
					NodeName: "node",
				},
			}),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			jobSpec := NewJobRunSpec(tc.args, tc.image, tc.pod)

			require.Equal(t, expJobName, jobSpec.Name)
			require.Equal(t, tc.args, jobSpec.Args)
			require.Equal(t, tc.image, jobSpec.Image)
			require.Equal(t, tc.pod.Node, jobSpec.Node)
			require.Equal(t, map[string]string{"job-name": expJobName}, jobSpec.Selectors)
		})
	}
}
