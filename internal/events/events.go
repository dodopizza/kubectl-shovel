package events

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

// Event is struct for events
type Event struct {
	Type    EventType `json:"type"`
	Message string    `json:"message"`
}

// EventType is type for further parsing
type EventType string

// available event types
const (
	Error     EventType = "error"
	Status    EventType = "status"
	Completed EventType = "completed"
)

// NewEvent is used to publish new event
func NewEvent(eventType EventType, message string) {
	data, _ := json.Marshal(Event{
		Type:    eventType,
		Message: message,
	})

	fmt.Println(string(data))
}

// NewStatusEvent is used to publish status event
func NewStatusEvent(message string) {
	NewEvent(Status, message)
}

// NewCompletedEvent is used to publish completed event
func NewCompletedEvent(message string) {
	NewEvent(Completed, message)
}

// NewErrorEvent is used to publish error event
func NewErrorEvent(err error, description string) {
	if description != "" {
		err = errors.Wrap(err, description)
	}
	NewEvent(Error, err.Error())
}

// GetEvent is used to read published event
func GetEvent(data string) (*Event, error) {
	event := &Event{}
	if err := json.Unmarshal([]byte(data), event); err != nil {
		return nil, errors.Wrap(err, "Can't parse event")
	}

	return event, nil
}
