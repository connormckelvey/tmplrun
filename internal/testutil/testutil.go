package testutil

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func MustReadFile(t *testing.T, path string) []byte {
	content, err := os.ReadFile(path)
	require.NoError(t, err)
	return content
}

func MustOpenFile(t *testing.T, path string) io.ReadCloser {
	content, err := os.Open(path)
	require.NoError(t, err)
	return content
}

func MustCloseFile(t *testing.T, f io.Closer) {
	require.NoError(t, f.Close())
}
