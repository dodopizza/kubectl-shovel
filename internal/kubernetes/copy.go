package kubernetes

import (
	"bytes"
	"fmt"
	"os"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	cp "k8s.io/kubectl/pkg/cmd/cp"
)

// Copy file from pod to local file
func (k8s *Client) Copy(podName, podFilePath, localFilePath string) error {
	ioStreams := genericclioptions.IOStreams{
		In:     &bytes.Buffer{},
		Out:    &bytes.Buffer{},
		ErrOut: os.Stdout,
	}
	opts := cp.NewCopyOptions(ioStreams)
	opts.Clientset = k8s.Clientset
	opts.ClientConfig = k8s.Config
	from := fmt.Sprintf("%s/%s:%s", k8s.Namespace, podName, podFilePath)

	return opts.Run([]string{from, localFilePath})
}
