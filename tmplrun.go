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

// TMPLRun represents a template runner that manages the rendering of templates.
type TMPLRun struct {
	fs fs.FS // fs is the filesystem from which templates will be loaded.
}

// New creates a new instance of TMPLRun with the given filesystem.
func New(fsys fs.FS) *TMPLRun {
	return &TMPLRun{fsys}
}

// RenderInput contains the input data required for rendering a template.
type RenderInput struct {
	Entrypoint string         // Entrypoint is the path to the main template file.
	Props      map[string]any // Props contains the properties to be passed to the template.
}

// Render renders the template specified by the input and writes the result to the given writer.
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
