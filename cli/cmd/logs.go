package cmd

import (
	"bufio"
	"fmt"
	"io"

	"github.com/dodopizza/kubectl-shovel/internal/events"
)

func handleLogs(readCloser io.ReadCloser) (string, error) {
	eventsChan := make(chan string)

	go func() {
		defer readCloser.Close()
		defer close(eventsChan)
		r := bufio.NewReader(readCloser)
		var err error
		for err != io.EOF {
			var event string
			event, err = r.ReadString('\n')
			if err != nil && err != io.EOF {
				return
			}

			eventsChan <- event
		}
	}()

	return processLogs(eventsChan)
}

func processLogs(eventsChan chan string) (string, error) {
	for rawEvent := range eventsChan {
		event, err := events.GetEvent(rawEvent)
		if err != nil {
			continue
		}

		switch event.Type {
		case events.Status:
			fmt.Println(event.Message)
		case events.Error:
			return "", fmt.Errorf("Error in job occurred: %s", event.Message)
		case events.Completed:
			return event.Message, nil
		default:
			return "", fmt.Errorf("Got unknown event type: %s", event.Type)
		}
	}

	return "", nil
}
