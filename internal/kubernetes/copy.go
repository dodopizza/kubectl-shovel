package kubernetes

import (
	"bytes"
	"fmt"
	"os"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/cmd/cp"
)

// CopyFromPod - copy file from pod to local file
func (k *Client) CopyFromPod(podName, podFilePath, localFilePath string) error {
	from := buildPodPath(k.Namespace, podName, podFilePath)

	return k.copy(from, localFilePath)
}

// CopyToPod - copy local file to pod
func (k *Client) CopyToPod(localFilePath, podName, podFilePath string) error {
	to := buildPodPath(k.Namespace, podName, podFilePath)

	return k.copy(localFilePath, to)
}

func (k *Client) copy(from, to string) error {
	ioStreams := genericclioptions.IOStreams{
		In:     &bytes.Buffer{},
		Out:    &bytes.Buffer{},
		ErrOut: os.Stdout,
	}
	opts := cp.NewCopyOptions(ioStreams)
	opts.Clientset = k.Clientset
	opts.ClientConfig = k.Config

	return opts.Run([]string{from, to})
}

func buildPodPath(namespace, podName, podFilePath string) string {
	return fmt.Sprintf("%s/%s:%s", namespace, podName, podFilePath)
}
