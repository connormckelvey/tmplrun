package tmplrun

import (
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/connormckelvey/tmplrun/internal/containers"
)

type hooks struct {
	tr        *TMPLRun
	fileStack *containers.Stack[string]
}

func (th *hooks) resolve(name string) string {
	// relative to importing file
	currentDir := filepath.Dir(th.fileStack.Peek())
	return filepath.Join(currentDir, name)
}

func (th *hooks) Include(name string) (string, error) {
	rel := th.resolve(name)
	b, err := fs.ReadFile(th.tr.fs, rel)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (th *hooks) Render(src string, props map[string]any) (string, error) {
	doc, err := th.tr.parse(strings.NewReader(src))
	if err != nil {
		return "", err
	}
	return th.tr.render(doc, props)
}
