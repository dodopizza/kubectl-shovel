package watchdog

import (
	"context"
	"errors"
	"github.com/dodopizza/kubectl-shovel/internal/globals"
	"io/ioutil"
	"os"
	"time"

	"github.com/dodopizza/kubectl-shovel/internal/kubernetes"
)

type Operator struct {
	k8s      *kubernetes.Client
	podName  string
	interval time.Duration
}

func NewOperator(k8s *kubernetes.Client, podName string) *Operator {
	return &Operator{
		k8s:      k8s,
		podName:  podName,
		interval: pingInterval,
	}
}

func (o *Operator) Run(ctx context.Context) error {
	successCh := make(chan struct{}, 1)

	go o.run(ctx, successCh)

	for {
		select {
		case <-successCh:
		case <-ctx.Done():
			return nil
		case <-time.After(deadAfterDuration):
			return errors.New("there were some issues to send ping to pod for a long time")
		}
	}
}

func (o *Operator) run(ctx context.Context, successCh chan<- struct{}) {
	ticker := time.NewTicker(o.interval)
	defer ticker.Stop()
	defer close(successCh)

	for {
		select {
		case <-ticker.C:
			if err := o.ping(); err != nil {
				continue
			}
			successCh <- struct{}{}
		case <-ctx.Done():
			return
		}
	}
}

// ping leave mark that operator is alive. Periodically update file at pod
func (o *Operator) ping() error {
	file, err := ioutil.TempFile("", globals.PluginName)
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())

	return o.k8s.CopyToPod(file.Name(), o.podName, pingFile)
}
