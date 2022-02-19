package watchdog

import (
	"context"
	"errors"
	"os"
	"time"
)

func Watch(ctx context.Context) error {
	pingCh := make(chan struct{}, 1)

	go watch(ctx, pingCh)

	for {
		select {
		case <-pingCh:
		case <-ctx.Done():
			return nil
		case <-time.After(deadAfterDuration):
			return errors.New("there were no signals from operator for a long time")
		}
	}
}

func watch(ctx context.Context, pingCh chan<- struct{}) {
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()
	defer close(pingCh)

	for {
		select {
		case <-ticker.C:
			isAlive := isOperatorAlive()
			if isAlive {
				pingCh <- struct{}{}
			}
		case <-ctx.Done():
			return
		}
	}
}

func isOperatorAlive() bool {
	file, err := os.Stat(pingFile)
	if err != nil {
		return false
	}

	return !time.Now().After(file.ModTime().Add(pingInterval))
}
