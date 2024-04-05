package tmplrun

import (
	"path/filepath"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestFS(t *testing.T) {
	tests := []string{
		"**/**/*.go",
	}

	for _, test := range tests {
		g, err := filepath.Glob(test)
		assert.NoError(t, err)
		spew.Dump(test, g)
	}
}
