package tmplrun

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/connormckelvey/tmplrun/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testTMPLRunSuccess struct {
	name           string
	unmarshalProps func([]byte, any) error
	includes       []string
}

func (s testTMPLRunSuccess) path(paths ...string) string {
	return filepath.Join("testdata/success", string(s.name), filepath.Join(paths...))
}

func (s testTMPLRunSuccess) entrypoint(t *testing.T) io.ReadCloser {
	return testutil.MustOpenFile(t, s.path("input/entrypoint"))
}

func (s testTMPLRunSuccess) props(t *testing.T) map[string]any {
	b := testutil.MustReadFile(t, s.path("input/props"))
	var m map[string]any
	err := s.unmarshalProps(b, &m)
	require.NoError(t, err)
	return m
}

func (s testTMPLRunSuccess) expected(t *testing.T) []byte {
	return testutil.MustReadFile(t, s.path("expected"))
}

func TestTMPLRun(t *testing.T) {
	tests := []testTMPLRunSuccess{
		{"example1", json.Unmarshal, []string{}},
	}

	for _, test := range tests {
		t.Run(string(test.name), func(t *testing.T) {
			entry, props, expected := test.entrypoint(t), test.props(t), test.expected(t)
			defer testutil.MustCloseFile(t, entry)

			fsys := os.DirFS(".")
			tmpl := New(fsys)

			result, err := tmpl.Render(&RenderInput{
				Entrypoint: test.path("input/entrypoint"),
				Props:      props,
			})
			assert.NoError(t, err)

			assert.Equal(t, string(expected), result)

		})
	}

}
