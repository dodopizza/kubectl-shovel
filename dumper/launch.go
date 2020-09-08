package main

import (
	"encoding/base64"
	"io/ioutil"

	"github.com/dodopizza/kubectl-shovel/events"
)

const (
	output = "/output.gcdump"
)

func launch(containerID string) error {
	events.NewEvent(
		events.Status,
		"Looking for container fs",
	)
	err := mapContainerTmp(containerID)
	if err != nil {
		return err
	}

	events.NewEvent(
		events.Status,
		"Starting dotnet-gcdump",
	)
	err = makeGcDump(1, output)
	if err != nil {
		return err
	}
	dumpContent, err := ioutil.ReadFile(output)
	if err != nil {
		return err
	}

	events.NewEvent(
		events.Result,
		base64.StdEncoding.EncodeToString(dumpContent),
	)

	return nil
}
