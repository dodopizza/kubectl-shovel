package version

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_GetVersion(t *testing.T) {
	require.Equal(t, "undefined", GetVersion())

	Version = "v1.2.0"
	require.Equal(t, "v1.2.0", GetVersion())
}
