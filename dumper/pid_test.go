package main

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	testPath = "./testdata"
)

func Test_findPID(t *testing.T) {
	pid, err := findPID(
		"658114371bb991d98dcaa576f4da91cfc6ae09e41cc440a4bc04b4b5eda45843",
		filepath.Join(
			testPath,
			"proc",
		),
	)
	require.NoError(t, err)

	require.Equal(t, 6543, pid)
}

func Test_findPID_Errors(t *testing.T) {
	_, err := findPID(
		"123",
		filepath.Join(
			testPath,
			"proc",
		),
	)
	require.Error(t, err)
}
