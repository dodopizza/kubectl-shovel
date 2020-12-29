package cmd

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"

	"github.com/dodopizza/kubectl-shovel/dumper/utils"
	"github.com/dodopizza/kubectl-shovel/pkg/events"
)

const (
	output = "/output"
)

func launch(executable string, args ...string) error {
	if containerRuntime == "docker" {
		events.NewEvent(
			events.Status,
			"Looking for container fs",
		)
		err := utils.MapDockerContainerTmp(containerID)
		if err != nil {
			return err
		}
	}
	events.NewEvent(
		events.Status,
		fmt.Sprintf(
			"Starting %s in job",
			executable,
		),
	)

	err := utils.ExecCommand(
		executable,
		args...,
	)
	if err != nil {
		return err
	}
	result, err := ioutil.ReadFile(output)
	if err != nil {
		return err
	}
	events.NewEvent(
		events.Result,
		base64.StdEncoding.EncodeToString(result),
	)

	return nil
}
