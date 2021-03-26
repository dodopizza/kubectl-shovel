// +build integration

package integration_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/dodopizza/kubectl-shovel/cli/cmd"
)

var (
	sampleAppName  = "sample-app"
	sampleAppImage = "mcr.microsoft.com/dotnet/core/samples:aspnetapp"
	namespace      = "default"
	cliPath        = "./bin/kubectl-shovel"
	dumperImage    = "kubectl-shovel/dumper-integration-tests"
	tempDirPattern = "*-kubectl-shovel"
)

func TestMain(m *testing.M) {
	k8s := newTestKubeClient()
	labels := map[string]string{
		"app": sampleAppName,
	}

	sampleAppPod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:   sampleAppName,
			Labels: labels,
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

	fmt.Println("Deploying dotnet sample app to cluster...")
	_, err := k8s.CoreV1().Pods(namespace).Create(context.TODO(), sampleAppPod, metav1.CreateOptions{})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Waiting app to start...")
	_, err = k8s.WaitPod(labels)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	exitCode := m.Run()
	_ = k8s.CoreV1().Pods(namespace).Delete(
		context.TODO(),
		sampleAppName,
		metav1.DeleteOptions{PropagationPolicy: &deletePolicy},
	)

	os.Exit(exitCode)
}

func Test_TraceSubcommand(t *testing.T) {
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
	}
	os.Args = append([]string{os.Args[0]}, args...)
	rootCmd := cmd.NewShovelCommand()
	require.NoError(t, rootCmd.Execute())

	file, err := os.Stat(outputFilename)
	require.NoError(t, err)
	require.NotEmpty(t, file.Size())
}

func Test_GCDumpSubcommand(t *testing.T) {
	dir, _ := ioutil.TempDir("", tempDirPattern)
	defer os.RemoveAll(dir)
	outputFilename := filepath.Join(dir, "gcdump-test")
	args := []string{
		"gcdump",
		"--pod-name",
		sampleAppName,
		"--output",
		outputFilename,
		"--image",
		dumperImage,
	}
	os.Args = append([]string{os.Args[0]}, args...)
	rootCmd := cmd.NewShovelCommand()
	require.NoError(t, rootCmd.Execute())

	file, err := os.Stat(outputFilename)
	require.NoError(t, err)
	require.NotEmpty(t, file.Size())
}
