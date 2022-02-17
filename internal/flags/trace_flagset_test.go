package flags

import (
	"strings"
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/require"
)

func Test_TestFlagSet(t *testing.T) {
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
				"--duration", "00:00:00:10",
			},
		},
		{
			name: "Override process ID",
			args: []string{
				"--process-id", "5",
			},
			expArgs: []string{
				"--process-id", "5",
				"--duration", "00:00:00:10",
			},
		},
		{
			name: "Override buffersize",
			args: []string{
				"--buffersize", "1024",
			},
			expArgs: []string{
				"--process-id", "1",
				"--buffersize", "1024",
				"--duration", "00:00:00:10",
			},
		},
		{
			name: "Override buffersize",
			args: []string{
				"--clreventlevel", "warning",
			},
			expArgs: []string{
				"--process-id", "1",
				"--clreventlevel", "warning",
				"--duration", "00:00:00:10",
			},
		},
		{
			name: "Override buffersize",
			args: []string{
				"--clrevents", "gc+gchandle",
			},
			expArgs: []string{
				"--process-id", "1",
				"--clrevents", "gc+gchandle",
				"--duration", "00:00:00:10",
			},
		},
		{
			name: "1 minute duration",
			args: []string{
				"--duration", "00:00:01:00",
			},
			expArgs: []string{
				"--process-id", "1",
				"--duration", "00:00:01:00",
			},
		},
		{
			name: "5 days duration and 1 minute",
			args: []string{
				"--duration", "05:00:01:00",
			},
			expArgs: []string{
				"--process-id", "1",
				"--duration", "05:00:01:00",
			},
		},
		{
			name: "1 minute duration with units",
			args: []string{
				"--duration", "1m",
			},
			expArgs: []string{
				"--process-id", "1",
				"--duration", "00:00:01:00",
			},
		},
		{
			name: "Duration with multiple units",
			args: []string{
				"--duration", "5h10s",
			},
			expArgs: []string{
				"--process-id", "1",
				"--duration", "00:05:00:10",
			},
		},
		{
			name: "Duration with ms units",
			args: []string{
				"--duration", "5m50ms",
			},
			expArgs: []string{
				"--process-id", "1",
				"--duration", "00:00:05:00",
			},
		},
		{
			name: "Override duration with units",
			args: []string{
				"--duration", "1m",
			},
			expArgs: []string{
				"--process-id", "1",
				"--duration", "00:00:01:00",
			},
		},
		{
			name: "Override format",
			args: []string{
				"--format", "Speedscope",
			},
			expArgs: []string{
				"--process-id", "1",
				"--duration", "00:00:00:10",
				"--format", "Speedscope",
			},
		},
		{
			name: "Override profile",
			args: []string{
				"--profile", "gc-verbose",
			},
			expArgs: []string{
				"--process-id", "1",
				"--duration", "00:00:00:10",
				"--profile", "gc-verbose",
			},
		},
		{
			name: "Override providers",
			args: []string{
				"--providers",
				strings.Join(
					[]string{
						"System.Runtime:0:1:EventCounterIntervalSec=1",
						"Microsoft-Windows-DotNETRuntime:0:1",
						"Microsoft-DotNETCore-SampleProfiler:0:1",
					},
					",",
				),
			},
			expArgs: []string{
				"--process-id", "1",
				"--duration", "00:00:00:10",
				"--providers",
				strings.Join(
					[]string{
						"System.Runtime:0:1:EventCounterIntervalSec=1",
						"Microsoft-Windows-DotNETRuntime:0:1",
						"Microsoft-DotNETCore-SampleProfiler:0:1",
					},
					",",
				),
			},
		},
		{
			name: "Override profile",
			args: []string{
				"--profile", "gc-verbose",
			},
			expArgs: []string{
				"--process-id", "1",
				"--duration", "00:00:00:10",
				"--profile", "gc-verbose",
			},
		},
		{
			name: "Set multiple flags",
			args: []string{
				"--buffersize", "1024",
				"--duration", "00:00:01:00",
				"--format", "Speedscope",
				"--profile", "gc-verbose",
			},
			expArgs: []string{
				"--process-id", "1",
				"--buffersize", "1024",
				"--duration", "00:00:01:00",
				"--format", "Speedscope",
				"--profile", "gc-verbose",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)
			trace := NewTraceFlagSet()
			flagSet.AddFlagSet(trace.GetFlags())

			require.NoError(t, flagSet.Parse(tc.args))
			require.Equal(t, tc.expArgs, trace.GetArgs())
		})
	}
}

func Test_TraceFlagSet_Errors(t *testing.T) {
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
			name: "Bad buffersize",
			args: []string{
				"--buffersize", "wrong",
			},
		},
		{
			name: "Bad CLR Event Level",
			args: []string{
				"--clreventlevel", "high",
			},
		},
		{
			name: "Empty CLR Event Level",
			args: []string{
				"--clreventlevel", "",
			},
		},
		{
			name: "Empty CLR Events",
			args: []string{
				"--clrevents", "",
			},
		},
		{
			name: "Wrong duration format",
			args: []string{
				"--duration", "00:00:00",
			},
		},
		{
			name: "Not numbers in duration",
			args: []string{
				"--duration", "00:00:00:aa",
			},
		},
		{
			name: "Just numbers in duration",
			args: []string{
				"--duration", "100",
			},
		},
		{
			name: "Too low duration",
			args: []string{
				"--duration", "100ms",
			},
		},
		{
			name: "Empty duration",
			args: []string{
				"--duration", "",
			},
		},
		{
			name: "Wrong format",
			args: []string{
				"--format", "Scopespeed",
			},
		},
		{
			name: "Empty format",
			args: []string{
				"--format", "",
			},
		},
		{
			name: "Wrong profile",
			args: []string{
				"--profile", "a4",
			},
		},
		{
			name: "Empty profile",
			args: []string{
				"--profile", "",
			},
		},
		{
			name: "Empty providers",
			args: []string{
				"--providers", "",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)
			trace := NewTraceFlagSet()
			flagSet.AddFlagSet(trace.GetFlags())

			require.Error(t, flagSet.Parse(tc.args))
		})
	}
}
