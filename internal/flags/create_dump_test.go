package flags

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/require"
)

func Test_CreateDumpFlagSetBinary(t *testing.T) {
	testCases := []struct {
		name          string
		args          []string
		expArgsTool   []string
		expArgsBinary []string
	}{
		{
			name:        "Defaults",
			args:        []string{},
			expArgsTool: []string{"--process-id", "1"},
			expArgsBinary: []string{
				"1",
			},
		},
		{
			name: "Override process ID",
			args: []string{
				"--process-id", "5",
			},
			expArgsTool: []string{"--process-id", "5"},
			expArgsBinary: []string{
				"5",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tool := NewCreateDump()
			flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)
			flagSet.AddFlagSet(tool.GetFlags())

			// require no error for parsing
			err := flagSet.Parse(tc.args)
			require.NoError(t, err)

			// format args for tool
			args := NewArgs()
			tool.FormatArgs(args, FormatArgsTypeTool)
			require.Equal(t, tc.expArgsTool, args.Get())

			// format args for binary
			args = NewArgs()
			tool.FormatArgs(args, FormatArgsTypeBinary)
			require.Equal(t, tc.expArgsBinary, args.Get())
		})
	}
}

func Test_CreateDumpFlagSetTool(t *testing.T) {
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
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			args := NewArgs()
			tool := NewCreateDump()
			flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)
			flagSet.AddFlagSet(tool.GetFlags())

			err := flagSet.Parse(tc.args)
			tool.FormatArgs(args, FormatArgsTypeTool)

			require.NoError(t, err)
			require.Equal(t, tc.expArgs, args.Get())
		})
	}
}

func Test_CreateDumpFlagSet_Errors(t *testing.T) {
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
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)
			createdump := NewCreateDump()
			flagSet.AddFlagSet(createdump.GetFlags())

			require.Error(t, flagSet.Parse(tc.args))
		})
	}
}
