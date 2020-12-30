package events

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var testCases = []struct {
	eventType   EventType
	message     string
	eventString string
}{
	{
		eventType:   Error,
		message:     "text error message",
		eventString: `{"type":"error","message":"text error message"}`,
	},
	{
		eventType:   Result,
		message:     "someresult",
		eventString: `{"type":"result","message":"someresult"}`,
	},
	{
		eventType:   "custom",
		message:     "anything",
		eventString: `{"type":"custom","message":"anything"}`,
	},
	{
		eventType:   "",
		message:     "",
		eventString: `{"type":"","message":""}`,
	},
}

func Test_NewEvent(t *testing.T) {
	for _, tc := range testCases {
		testStdout, writer, err := os.Pipe()
		require.NoError(t, err)
		osStdout := os.Stdout
		os.Stdout = writer
		defer func() {
			os.Stdout = osStdout
		}()

		NewEvent(tc.eventType, tc.message)

		writer.Close()

		var buf bytes.Buffer
		io.Copy(&buf, testStdout)
		actual := buf.String()
		require.Equal(t, tc.eventString+"\n", actual)
	}
}

func Test_GetEvent(t *testing.T) {
	for _, tc := range testCases {
		actual, err := GetEvent(tc.eventString)
		require.NoError(t, err)
		require.Equal(
			t,
			&Event{
				Type:    tc.eventType,
				Message: tc.message,
			},
			actual,
		)
	}
}

func Test_GetEvent_Error(t *testing.T) {
	errorCases := []string{
		"{",
		"not a json",
		`{"broken":"json"`,
	}
	for _, tc := range errorCases {
		_, err := GetEvent(tc)
		require.Error(t, err)
	}
}
