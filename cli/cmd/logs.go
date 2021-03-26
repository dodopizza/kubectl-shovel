package cmd

import (
	"bufio"
	"fmt"
	"io"

	"github.com/dodopizza/kubectl-shovel/internal/events"
	"github.com/pkg/errors"
)

func handleLogs(readCloser io.ReadCloser) (string, error) {
	eventsChan := make(chan string)
	defer close(eventsChan)

	go func() {
		defer readCloser.Close()
		r := bufio.NewReader(readCloser)
		for {
			bytes, err := r.ReadBytes('\n')
			if err != nil {
				return
			}

			eventsChan <- string(bytes)
		}
	}()

	return processLogs(eventsChan)
}

func processLogs(eventsChan chan string) (string, error) {
	var resultFilePath string
LOOP:
	for rawEvent := range eventsChan {
		event, err := events.GetEvent(rawEvent)
		if err != nil {
			return "", errors.Wrap(err, "Got malformed event")
		}

		switch event.Type {
		case events.Status:
			fmt.Println(event.Message)
		case events.Error:
			return "", fmt.Errorf("Error in job occurred: %s", event.Message)
		case events.Completed:
			resultFilePath = event.Message
			fmt.Printf("Results located at %s", resultFilePath)
			break LOOP
		default:
			return "", fmt.Errorf("Got unknown event type: %s", event.Type)
		}
	}

	return resultFilePath, nil
}
