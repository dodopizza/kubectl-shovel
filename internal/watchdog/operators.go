package watchdog

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/dodopizza/kubectl-shovel/internal/globals"
	"github.com/dodopizza/kubectl-shovel/internal/kubernetes"
)

func NewWatcher() *Operator {
	return NewOperator(
		func() bool {
			file, err := os.Stat(pingFile)
			if err != nil {
				return false
			}

			return !time.Now().After(file.ModTime().Add(pingInterval))
		},
		deadline,
		checkInterval,
	)
}

func NewPinger(k8s *kubernetes.Client, pod string) *Operator {
	return NewOperator(
		func() bool {
			file, err := ioutil.TempFile("", globals.PluginName)
			if err != nil {
				return false
			}
			defer os.Remove(file.Name())

			err = k8s.CopyToPod(file.Name(), pod, pingFile)
			return err != nil
		},
		deadline,
		pingInterval,
	)
}
