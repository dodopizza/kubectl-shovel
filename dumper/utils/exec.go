package utils

import (
	"bytes"
	"os/exec"

	"github.com/pkg/errors"
)

func ExecCommand(executable string, args ...string) error {
	cmd := exec.Command(
		executable,
		args...,
	)

	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stdout

	err := cmd.Run()
	if err != nil {
		return errors.Wrap(err, stdout.String())
	}

	return nil
}
