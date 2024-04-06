package tmplrun

import (
	"bytes"
	"io/fs"
	"path/filepath"
)

type hooks struct {
	tr          *TMPLRun
	currentFile string
}

func (th *hooks) resolve(name string) string {
	currentDir := filepath.Dir(th.currentFile)
	return filepath.Join(currentDir, name)
}

// Include resolves and includes the template specified by name.
func (th *hooks) Include(name string) (string, error) {
	rel := th.resolve(name)
	b, err := fs.ReadFile(th.tr.fs, rel)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Render resolves, parses, and renders the template specified by name with the given properties.
func (th *hooks) Render(name string, props map[string]any) (string, error) {
	rel := th.resolve(name)
	src, err := fs.ReadFile(th.tr.fs, rel)
	if err != nil {
		return "", err
	}

	doc, err := th.tr.parse(bytes.NewReader(src))
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = th.tr.render(&buf, rel, doc, props)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
