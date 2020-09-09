package main

import (
	"flag"

	"github.com/dodopizza/kubectl-shovel/events"
)

func main() {
	var containerID string
	var tool string
	flag.StringVar(&containerID, "container-id", containerID, "ContainerID for creating dump")
	flag.StringVar(&tool, "tool", tool, "ContainerID for creating dump")
	flag.Parse()
	if containerID == "" || tool == "" {
		events.NewEvent(events.Error, "ContainerID or tool not defined")
	}

	err := launch(containerID, tool)
	if err != nil {
		events.NewEvent(events.Error, err.Error())
	}
}
