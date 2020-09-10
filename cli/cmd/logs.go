package cmd

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/dodopizza/kubectl-shovel/pkg/events"
)

func handleLogs(readCloser io.ReadCloser, output string) {
	eventsChan := make(chan string)
	done := make(chan struct{})

	go processLogs(eventsChan, done, output)
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

	<-done
	close(eventsChan)
}

func processLogs(eventsChan chan string, done chan struct{}, output string) {
	for rawEvent := range eventsChan {
		event, err := events.GetEvent(rawEvent)
		if err != nil {
			continue
		}
		switch event.Type {
		case events.Status:
			fmt.Println(event.Message)
		case events.Error:
			fmt.Println("Error occurred ", event.Message)
			done <- struct{}{}
		case events.Result:
			data, err := base64.StdEncoding.DecodeString(event.Message)
			if err != nil {
				fmt.Printf("Failed while decoding dump: %v\n", err)
			}

			err = ioutil.WriteFile(
				output,
				data,
				0777,
			)
			if err != nil {
				fmt.Printf("Failed whil writing to file: %v\n", err)
			}
			done <- struct{}{}
		}
	}
}
