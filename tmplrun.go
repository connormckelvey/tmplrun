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

func (tr *TMPLRun) Render(input *RenderInput) (string, error) {
	f, err := tr.fs.Open(input.Entrypoint)
	if err != nil {
		return "", err
	}
	defer f.Close()

	doc, err := tr.parse(f)
	if err != nil {
		return "", err
	}
	result, err := tr.render(input.Entrypoint, doc, input.Props)
	if err != nil {
		return "", err
	}
	return result, err
}

func (tr *TMPLRun) Run(reader io.Reader, props map[string]any) (string, error) {
	doc, err := tr.parse(reader)
	if err != nil {
		return "", err
	}
	return tr.render(".", doc, props)
}

func (tr *TMPLRun) parse(r io.Reader) (*ast.Document, error) {
	lex := lexer.New(r)
	par := parser.New(lex)
	return par.Parse()
}

func (tr *TMPLRun) render(currentFile string, doc *ast.Document, props map[string]any) (string, error) {
	hooks := &hooks{
		tr:          tr,
		currentFile: currentFile,
	}
	return evaluator.
		New(driver.NewGoja(), hooks).
		Render(doc, evaluator.NewEnvironment(tr.fs, props, hooks))
}
