package main

import (
	"flag"
	"log"
)

func main() {
	containerID := flag.String("containerID", "", "conrainerID for creating dump")

	if *containerID == "" {
		log.Fatal("containerID is empty")
	}

	launch(*containerID)
}
