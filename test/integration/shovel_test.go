// +build integration

package integration_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/dodopizza/kubectl-shovel/dumper/utils"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	sampleAppName  = "sample-app"
	sampleAppImage = "mcr.microsoft.com/dotnet/core/samples:aspnetapp"
	namespace      = "default"
	cliPath        = "./bin/kubectl-shovel"
	dumperImage    = "kubectl-shovel/integration-tests"
	tempDirPattern = "*-kubectl-shovel"
)

func TestMain(m *testing.M) {
	k8s := newTestKubeClient()

	sampleAppPod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: sampleAppName,
			Labels: map[string]string{
				"app": sampleAppName,
			},
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  sampleAppName,
					Image: sampleAppImage,
				},
			},
		},
	}

	_, err := k8s.CoreV1().Pods(namespace).Create(context.TODO(), sampleAppPod, metav1.CreateOptions{})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	<-time.After(10 * time.Second)
	exitCode := m.Run()
	_ = k8s.CoreV1().Pods(namespace).Delete(context.TODO(), sampleAppName, metav1.DeleteOptions{PropagationPolicy: &deletePolicy})

	os.Exit(exitCode)
}

func Test_TraceSubcommand(t *testing.T) {
	dir, _ := ioutil.TempDir("", tempDirPattern)
	defer os.RemoveAll(dir)
	outputFilename := filepath.Join(dir, "trace-test")
	err := utils.ExecCommand(
		cliPath,
		"trace",
		"--pod-name",
		sampleAppName,
		"--output",
		outputFilename,
		"--image",
		dumperImage,
	)

	require.NoError(t, err)

	file, err := os.Stat(outputFilename)
	require.NoError(t, err)
	require.NotEmpty(t, file.Size())
}

func Test_GCDumpSubcommand(t *testing.T) {
	dir, _ := ioutil.TempDir("", tempDirPattern)
	defer os.RemoveAll(dir)
	outputFilename := filepath.Join(dir, "gcdump-test")
	err := utils.ExecCommand(
		cliPath,
		"gcdump",
		"--pod-name",
		sampleAppName,
		"--output",
		outputFilename,
		"--image",
		dumperImage,
	)

	require.NoError(t, err)

	file, err := os.Stat(outputFilename)
	require.NoError(t, err)
	require.NotEmpty(t, file.Size())
}
