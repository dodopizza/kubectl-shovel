package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	var containerID string
	flag.StringVar(&containerID, "container-id", containerID, "ContainerID for creating dump")
	flag.Parse()

	if containerID == "" {
		log.Fatal("ContainerID is empty")
	}

	fmt.Println(launch(containerID))
}
