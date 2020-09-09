package main

import (
	"bytes"
	"os/exec"

	"github.com/dodopizza/kubectl-shovel/events"
)

func runCommand(executable string, args ...string) error {
	cmd := exec.Command(
		executable,
		args...,
	)

	var Stdout bytes.Buffer
	cmd.Stdout = &Stdout

	err := cmd.Run()

	events.NewEvent(
		events.Status,
		Stdout.String(),
	)

	if err != nil {
		return err
	}

	return nil
}
