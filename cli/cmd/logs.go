package cmd

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/dodopizza/kubectl-shovel/internal/events"
	"github.com/pkg/errors"
)

func handleLogs(readCloser io.ReadCloser, output string) error {
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

	return processLogs(eventsChan, output)
}

func processLogs(eventsChan chan string, output string) error {
LOOP:
	for rawEvent := range eventsChan {
		event, err := events.GetEvent(rawEvent)
		if err != nil {
			return errors.Wrap(err, "Got malformed event")
		}

		switch event.Type {
		case events.Status:
			fmt.Println(event.Message)
		case events.Error:
			return fmt.Errorf("Error in job occurred: %s", event.Message)
		case events.Result:
			if err := saveResult(event.Message, output); err != nil {
				return errors.Wrap(err, "Error occurred while saving results")
			}

			fmt.Printf("Result successfully written to %s\n", output)
			break LOOP
		default:
			return fmt.Errorf("Got unknown event type: %s", event.Type)
		}
	}

	return nil
}

func saveResult(message, output string) error {
	data, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		return errors.Wrap(err, "Failed while decoding dump")
	}

	if err := ioutil.WriteFile(
		output,
		data,
		0777,
	); err != nil {
		return errors.Wrap(err, "Failed while writing to file")
	}

	return nil
}
