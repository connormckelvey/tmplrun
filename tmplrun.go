package tmplrun

import (
	"io"
	"io/fs"

	"github.com/connormckelvey/tmplrun/ast"
	"github.com/connormckelvey/tmplrun/evaluator"
	"github.com/connormckelvey/tmplrun/evaluator/driver"
	"github.com/connormckelvey/tmplrun/lexer"
	"github.com/connormckelvey/tmplrun/parser"
)

type TMPLRun struct {
	fs fs.FS
}

func New(fsys fs.FS) *TMPLRun {
	return &TMPLRun{fsys}
}

type RenderInput struct {
	Entrypoint string
	Props      map[string]any
}

func (tr *TMPLRun) Render(w io.Writer, input *RenderInput) error {
	f, err := tr.fs.Open(input.Entrypoint)
	if err != nil {
		return err
	}
	defer f.Close()

	doc, err := tr.parse(f)
	if err != nil {
		return err
	}
	err = tr.render(w, input.Entrypoint, doc, input.Props)
	if err != nil {
		return err
	}
	return nil
}

func (tr *TMPLRun) parse(r io.Reader) (*ast.Document, error) {
	lex := lexer.New(r)
	par := parser.New(lex)
	return par.Parse()
}

func (tr *TMPLRun) render(w io.Writer, currentFile string, doc *ast.Document, props map[string]any) error {
	hooks := &hooks{
		tr:          tr,
		currentFile: currentFile,
	}
	ev := evaluator.New(driver.NewGoja(), hooks)
	res, err := ev.Render(doc, evaluator.NewEnvironment(tr.fs, props, hooks))
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(res))
	if err != nil {
		return err
	}

	return nil
}
