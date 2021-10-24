package flags

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/require"
)

func Test_ManagedDumpFlagSet(t *testing.T) {
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
				"--type", "Full",
			},
		},
		{
			name: "Override process ID",
			args: []string{
				"--process-id", "5",
			},
			expArgs: []string{
				"--process-id", "5",
				"--type", "Full",
			},
		},
		{
			name: "Override Type",
			args: []string{
				"--type", "Heap",
			},
			expArgs: []string{
				"--process-id", "1",
				"--type", "Heap",
			},
		},
		{
			name: "Override Diagnostics",
			args: []string{
				"--diag",
			},
			expArgs: []string{
				"--process-id", "1",
				"--diag",
				"--type", "Full",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)
			gc := NewManagedDumpFlagSet()
			flagSet.AddFlagSet(gc.Parse())

			require.NoError(t, flagSet.Parse(tc.args))
			require.Equal(t, tc.expArgs, gc.Args())
		})
	}
}

func Test_ManagedDumpFlagSet_Errors(t *testing.T) {
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
			name: "Bad Type",
			args: []string{
				"--type", "invalid",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)
			gc := NewGCDumpFlagSet()
			flagSet.AddFlagSet(gc.Parse())

			require.Error(t, flagSet.Parse(tc.args))
		})
	}
}