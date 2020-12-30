package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ExecCommand(t *testing.T) {
	err := ExecCommand(
		"/bin/bash",
		"-c",
		"exit 0",
	)

	require.NoError(t, err)
}

func Test_ExecCommand_Error(t *testing.T) {
	stdoutText := "Exec failure text in stdout"
	stderrText := "Exec failure text in stderr"
	err := ExecCommand(
		"/bin/bash",
		"-c",
		fmt.Sprintf(
			"echo '%s' && echo '%s' >&2 && exit 1",
			stdoutText,
			stderrText,
		),
	)

	require.Error(t, err)
	require.Contains(
		t,
		err.Error(),
		stdoutText,
	)
	require.Contains(
		t,
		err.Error(),
		stderrText,
	)
}
