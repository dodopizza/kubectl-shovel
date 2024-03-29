package flags

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/require"
)

func Test_DumpFlagSet(t *testing.T) {
	testCases := []struct {
		name    string
		args    []string
		expArgs []string
	}{
		{
			name:    "Defaults",
			args:    []string{},
			expArgs: []string{"--process-id", "1", "--type", "Full"},
		},
		{
			name:    "Override process ID",
			args:    []string{"--process-id", "5"},
			expArgs: []string{"--process-id", "5", "--type", "Full"},
		},
		{
			name:    "Override Type",
			args:    []string{"--type", "Heap"},
			expArgs: []string{"--process-id", "1", "--type", "Heap"},
		},
		{
			name:    "Override Diagnostics",
			args:    []string{"--diag"},
			expArgs: []string{"--process-id", "1", "--diag", "--type", "Full"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tool := NewDotnetDump()
			flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)
			flagSet.AddFlagSet(tool.GetFlags())

			// require no error for parsing
			err := flagSet.Parse(tc.args)
			require.NoError(t, err)

			// format args for tool
			args := NewArgs()
			tool.FormatArgs(args, FormatArgsTypeTool)
			require.Equal(t, tc.expArgs, args.Get())

			// format args for binary
			args = NewArgs()
			tool.FormatArgs(args, FormatArgsTypeBinary)
			require.Equal(t, append([]string{"collect"}, tc.expArgs...), args.Get())
		})
	}
}

func Test_DumpFlagSet_Errors(t *testing.T) {
	testCases := []struct {
		name string
		args []string
	}{
		{
			name: "Bad process ID",
			args: []string{"--process-id", "a"},
		},
		{
			name: "Empty process ID",
			args: []string{"--process-id", ""},
		},
		{
			name: "Bad Type",
			args: []string{"--type", "invalid"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tool := NewDotnetDump()
			flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)
			flagSet.AddFlagSet(tool.GetFlags())

			require.Error(t, flagSet.Parse(tc.args))
		})
	}
}
