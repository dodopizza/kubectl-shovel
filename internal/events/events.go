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
	Error  EventType = "error"
	Status EventType = "status"
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

// GetEvent is used to read published event
func GetEvent(data string) (*Event, error) {
	event := &Event{}
	if err := json.Unmarshal([]byte(data), event); err != nil {
		return nil, errors.Wrap(err, "Can't parse event")
	}

	return event, nil
}
