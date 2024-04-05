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
	// relative to importing file
	currentDir := filepath.Dir(th.currentFile)
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
	return th.tr.render(rel, doc, props)
}
