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
	events.NewEvent(
		events.Status,
		"Looking for and mapping container fs",
	)
	if err := mapContainerTmp(containerInfo); err != nil {
		return err
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
	if err := utils.ExecCommand(
		executable,
		args...,
	); err != nil {
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
