package main

import (
	"flag"

	"github.com/dodopizza/kubectl-shovel/events"
)

func main() {
	var containerID string
	flag.StringVar(&containerID, "container-id", containerID, "ContainerID for creating dump")
	flag.Parse()

	if containerID == "" {
		events.NewEvent(events.Error, "ContainerID is empty")
	}

	err := launch(containerID)
	if err != nil {
		events.NewEvent(events.Error, err.Error())
	}
}
