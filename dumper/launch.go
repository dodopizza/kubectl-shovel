package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"

	"github.com/dodopizza/kubectl-shovel/events"
)

const ()

func launch(containerID, tool string) error {
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
		fmt.Sprintf(
			"Starting %s",
			tool,
		),
	)

	var output string
	switch tool {
	case dotnetGCDumpBinary:
		output = "/output.gcdump"
		err = makeGCDump(1, output)
	case dotnetTraceBinary:
		output = "/output.nettrace"
		err = makeTrace(1, output)
	}
	if err != nil {
		return err
	}
	content, err := ioutil.ReadFile(output)
	if err != nil {
		return err
	}

	events.NewEvent(
		events.Result,
		base64.StdEncoding.EncodeToString(content),
	)

	return nil
}
