package flags

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/require"
)

func Test_CreateDumpFlagSet(t *testing.T) {
	testCases := []struct {
		name    string
		args    []string
		expArgs []string
	}{
		{
			name: "Defaults",
			args: []string{},
			expArgs: []string{
				"1",
			},
		},
		{
			name: "Override process ID",
			args: []string{
				"--process-id", "5",
			},
			expArgs: []string{
				"5",
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
