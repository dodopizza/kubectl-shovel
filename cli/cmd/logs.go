package cmd

import (
	"bufio"
	"fmt"
	"io"

	"github.com/dodopizza/kubectl-shovel/internal/events"
)

func handleLogs(rc io.ReadCloser) (string, error) {
	entries := make(chan string, 1)

	go func() {
		defer rc.Close()
		defer close(entries)
		r := bufio.NewReader(rc)

		for {
			payload, err := r.ReadString('\n')

			if err == io.EOF {
				break
			}
			if err != nil {
				return
			}

			entries <- payload
		}
	}()

	return processLogs(entries)
}

func processLogs(entries chan string) (string, error) {
	for entry := range entries {
		event, err := events.GetEvent(entry)
		if err != nil {
			continue
		}

		switch event.Type {
		case events.Status:
			fmt.Println(event.Message)
		case events.Error:
			return "", fmt.Errorf("error in job occurred: %s", event.Message)
		case events.Completed:
			return event.Message, nil
		default:
			return "", fmt.Errorf("got unknown event type: %s", event.Type)
		}
	}

	return "", nil
}
