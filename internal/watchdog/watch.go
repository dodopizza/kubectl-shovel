package watchdog

import (
	"errors"
	"os"
	"time"
)

func Watch() error {
	pingCh := make(chan struct{}, 1)
	go watch(pingCh)
	for {
		select {
		case <-pingCh:
		case <-time.After(deadAfterDuration):
			return errors.New("There were no signals from operator for a long time")
		}
	}
}

func watch(pingCh chan<- struct{}) {
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()
	for range ticker.C {
		isAlive := isOperatorAlive()
		if isAlive {
			pingCh <- struct{}{}
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
