package kubernetes

import (
	"bytes"
	"fmt"
	"os"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	cp "k8s.io/kubectl/pkg/cmd/cp"
)

// CopyFromPod - copy file from pod to local file
func (k8s *Client) CopyFromPod(podName, podFilePath, localFilePath string) error {
	from := buildPodPath(k8s.Namespace, podName, podFilePath)

	return k8s.copy(from, localFilePath)
}

// CopyToPod - copy local file to pod
func (k8s *Client) CopyToPod(localFilePath, podName, podFilePath string) error {
	to := buildPodPath(k8s.Namespace, podName, podFilePath)

	return k8s.copy(localFilePath, to)
}

func buildPodPath(namespace, podName, podFilePath string) string {
	return fmt.Sprintf("%s/%s:%s", namespace, podName, podFilePath)
}

func (k8s *Client) copy(from, to string) error {
	ioStreams := genericclioptions.IOStreams{
		In:     &bytes.Buffer{},
		Out:    &bytes.Buffer{},
		ErrOut: os.Stdout,
	}
	opts := cp.NewCopyOptions(ioStreams)
	opts.Clientset = k8s.Clientset
	opts.ClientConfig = k8s.Config

	return opts.Run([]string{from, to})
}
