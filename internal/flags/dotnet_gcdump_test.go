package flags

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/require"
)

func Test_GCDumpFlagSet(t *testing.T) {
	testCases := []struct {
		name    string
		args    []string
		expArgs []string
	}{
		{
			name: "Defaults",
			args: []string{},
			expArgs: []string{
				"--process-id", "1",
			},
		},
		{
			name: "Override process ID",
			args: []string{
				"--process-id", "5",
			},
			expArgs: []string{
				"--process-id", "5",
			},
		},
		{
			name: "Override timeout",
			args: []string{
				"--timeout", "120",
			},
			expArgs: []string{
				"--process-id", "1",
				"--timeout", "120",
			},
		},
		{
			name: "Override timeout in seconds",
			args: []string{
				"--timeout", "120s",
			},
			expArgs: []string{
				"--process-id", "1",
				"--timeout", "120",
			},
		},
		{
			name: "Override timeout in minutes",
			args: []string{
				"--timeout", "2m",
			},
			expArgs: []string{
				"--process-id", "1",
				"--timeout", "120",
			},
		},
		{
			name: "Override timeout with milliseconds",
			args: []string{
				"--timeout", "60s50ms",
			},
			expArgs: []string{
				"--process-id", "1",
				"--timeout", "60",
			},
		},
		{
			name: "Override timeout and process id",
			args: []string{
				"--process-id", "5",
				"--timeout", "10m",
			},
			expArgs: []string{
				"--process-id", "5",
				"--timeout", "600",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)
			gc := NewDotnetGCDump()
			flagSet.AddFlagSet(gc.GetFlags())

			require.NoError(t, flagSet.Parse(tc.args))
			require.Equal(t, tc.expArgs, gc.GetArgs())
		})
	}
}

func Test_GCDumpFlagSet_Errors(t *testing.T) {
	testCases := []struct {
		name string
		args []string
	}{
		{
			name: "Bad process ID",
			args: []string{
				"--process-id", "a",
			},
		},
		{
			name: "Empty process ID",
			args: []string{
				"--process-id", "",
			},
		},
		{
			name: "Bad timeout",
			args: []string{
				"--timeout", "abc",
			},
		},
		{
			name: "Low timeout",
			args: []string{
				"--timeout", "5ms",
			},
		},
		{
			name: "Empty timeout",
			args: []string{
				"--timeout", "",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)
			gc := NewDotnetGCDump()
			flagSet.AddFlagSet(gc.GetFlags())

			require.Error(t, flagSet.Parse(tc.args))
		})
	}
}
