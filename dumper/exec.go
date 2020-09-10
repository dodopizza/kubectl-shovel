package main

import (
	"bytes"
	"os/exec"

	"github.com/pkg/errors"
)

func runCommand(executable string, args ...string) error {
	cmd := exec.Command(
		executable,
		args...,
	)

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		return errors.Wrap(err, stdout.String())
	}

	return nil
}
