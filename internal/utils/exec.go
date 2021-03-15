package utils

import (
	"bytes"
	"os/exec"

	"github.com/pkg/errors"
)

// ExecCommand is wrapper for running external processes
func ExecCommand(executable string, args ...string) error {
	cmd := exec.Command(
		executable,
		args...,
	)

	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stdout

	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, stdout.String())
	}

	return nil
}
