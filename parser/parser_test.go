package parser

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/connormckelvey/tmplrun/internal/testutil"
	"github.com/connormckelvey/tmplrun/lexer"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testParserSuccess string

func (s testParserSuccess) path(paths ...string) string {
	return filepath.Join("./testdata/success", string(s), filepath.Join(paths...))
}

func (s testParserSuccess) input(t *testing.T) []byte {
	return testutil.MustReadFile(t, s.path("input"))
}

func (s testParserSuccess) expected(t *testing.T) []byte {
	return testutil.MustReadFile(t, s.path("expected"))
}

func TestParserSuccess(t *testing.T) {
	tests := []testParserSuccess{
		"example1",
		"tags1",
		"nested1",
	}

	for _, test := range tests {
		t.Run(string(test), func(t *testing.T) {
			input, expected := test.input(t), test.expected(t)

			p := New(lexer.New(bytes.NewReader(input)))
			ast, err := p.Parse()
			assert.NoError(t, err)

			actual, err := json.Marshal(ast)
			require.NoError(t, err)

			assert.JSONEq(t, string(expected), string(actual))
		})
	}
}
