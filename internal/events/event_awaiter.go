package events

import (
	"bufio"
	"fmt"
	"io"
)

type EventAwaiter struct {
	entries chan string
}

func NewEventAwaiter() *EventAwaiter {
	return &EventAwaiter{
		entries: make(chan string),
	}
}

func (ep *EventAwaiter) AwaitCompletedEvent(stream io.Reader) (string, error) {
	go ep.read(stream)
	return ep.parse()
}

func (ep *EventAwaiter) read(stream io.Reader) {
	defer close(ep.entries)

	reader := bufio.NewReader(stream)

	for {
		payload, err := reader.ReadString('\n')

		if err == io.EOF {
			break
		}

		if err != nil {
			return
		}

		ep.entries <- payload
	}
}

func (ep *EventAwaiter) parse() (string, error) {
	for entry := range ep.entries {
		event, err := GetEvent(entry)
		if err != nil {
			continue
		}

		switch event.Type {
		case Status:
			fmt.Println(event.Message)
		case Error:
			return "", fmt.Errorf("error in job occurred: %s", event.Message)
		case Completed:
			return event.Message, nil
		default:
			return "", fmt.Errorf("got unknown event type: %s", event.Type)
		}
	}

	return "", nil
}
