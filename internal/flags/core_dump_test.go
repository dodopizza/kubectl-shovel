package flags

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/require"
)

func Test_CoreDumpFlagSet(t *testing.T) {
	testCases := []struct {
		name          string
		args          []string
		expArgsTool   []string
		expArgsBinary []string
	}{
		{
			name:          "Defaults",
			args:          []string{},
			expArgsTool:   []string{"--process-id", "1", "--type", "Full"},
			expArgsBinary: []string{"1", "--full"},
		},
		{
			name:          "Override process ID",
			args:          []string{"--process-id", "5"},
			expArgsTool:   []string{"--process-id", "5", "--type", "Full"},
			expArgsBinary: []string{"5", "--full"},
		},
		{
			name:          "Override dump type",
			args:          []string{"--type", "Triage"},
			expArgsTool:   []string{"--process-id", "1", "--type", "Triage"},
			expArgsBinary: []string{"1", "--triage"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tool := NewCoreDump()
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

func Test_CoreDumpFlagSet_Errors(t *testing.T) {
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
			name: "Bad type",
			args: []string{"--type", "invalid"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tool := NewCoreDump()
			flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)
			flagSet.AddFlagSet(tool.GetFlags())

			require.Error(t, flagSet.Parse(tc.args))
		})
	}
}
