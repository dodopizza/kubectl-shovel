package deadman

import (
	"io/ioutil"
	"os"

	"github.com/dodopizza/kubectl-shovel/internal/kubernetes"
)

var aliveFile = "/alive.tmp"

// Alive - leave mark, operator is alive. Periodically update alive file at pod.
func Alive(k8s *kubernetes.Client, podName string) error {
	file, err := ioutil.TempFile("", "kubectl-shovel")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())

	return k8s.CopyToPod(file.Name(), podName, aliveFile)
}
