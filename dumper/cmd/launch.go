package cmd

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/dodopizza/kubectl-shovel/internal/events"
	"github.com/dodopizza/kubectl-shovel/internal/utils"
)

func launch(executable string, args ...string) error {
	if containerRuntime == "docker" {
		events.NewEvent(
			events.Status,
			"Looking for container fs",
		)
		err := mapDockerContainerTmp(containerID)
		if err != nil {
			return err
		}
	}
	events.NewEvent(
		events.Status,
		fmt.Sprintf(
			"Running command: %s %s",
			executable,
			strings.Join(args, " "),
		),
	)

	// if we do not set proper file extension dotnet tools will do it anyway
	outputExtension := strings.TrimPrefix(executable, "dotnet-")
	output := fmt.Sprintf("/output.%s", outputExtension)
	args = append(
		args,
		"--output",
		output,
	)
	err := utils.ExecCommand(
		executable,
		args...,
	)
	if err != nil {
		return err
	}
	events.NewEvent(
		events.Status,
		"Gathering completed. Getting results",
	)
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
