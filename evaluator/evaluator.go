package evaluator

import (
	"strings"

	"github.com/connormckelvey/tmplrun/ast"
)

type Evaluator struct {
	driver Driver
	hooks  HooksAPI
}

func New(driver Driver, hooks HooksAPI) *Evaluator {
	return &Evaluator{
		driver: driver,
		hooks:  hooks,
	}
}

func (r *Evaluator) Render(doc *ast.Document, env *Environment) (string, error) {
	return r.renderChildren(doc, env)
}

func (r *Evaluator) renderChildren(n ast.Node, env *Environment) (string, error) {
	var code strings.Builder
	for _, child := range n.Children() {
		switch n := child.(type) {
		case *ast.TextNode:
			code.WriteString(n.Token.Literal)
		case *ast.ErrorNode:
			code.WriteString(n.Token.Literal)
		case *ast.TemplateNode:
			result, err := r.renderTemplateNode(n, env)
			if err != nil {
				return "", err
			}
			code.WriteString(result)
		}
	}
	return code.String(), nil
}

func (r *Evaluator) renderTemplateNode(n *ast.TemplateNode, env *Environment) (string, error) {
	code, err := r.renderChildren(n, env)
	if err != nil {
		return "", err
	}
	runtime, err := r.driver.CreateContext(env)
	if err != nil {
		return "", err
	}
	result, err := runtime.Eval(code)
	if err != nil {
		return "", err
	}
	if env.hasErrors() {
		return "", env.errors[0]
	}
	return result, nil
}
