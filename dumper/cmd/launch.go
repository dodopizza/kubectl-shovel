package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/dodopizza/kubectl-shovel/internal/events"
	"github.com/dodopizza/kubectl-shovel/internal/utils"
	"github.com/dodopizza/kubectl-shovel/internal/watchdog"
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
		"Gathering completed",
	)
	_, err := ioutil.ReadFile(output)
	if err != nil {
		return err
	}
	events.NewEvent(
		events.Completed,
		output,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := watchdog.Watch(ctx); err != nil {
		events.NewEvent(
			events.Error,
			err.Error(),
		)
		return err
	}

	return nil
}
