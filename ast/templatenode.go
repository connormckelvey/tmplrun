package ast

import (
	"encoding/json"

	"github.com/connormckelvey/tmplrun/lexer"
	"github.com/connormckelvey/tmplrun/token"
)

type TemplateNode struct {
	Token    *token.Token
	Closed   bool
	Ident    string
	children []Node
}

func parseTemplateTag(lit string) (ok bool, ident string, padding string) {
	match := lexer.OpenTagPattern.FindStringSubmatch(lit)
	if match == nil {
		return false, "", ""
	}
	return true, match[2], match[3]
}

func NewTemplateNode(tok *token.Token) *TemplateNode {
	_, ident, _ := parseTemplateTag(tok.Literal)
	return &TemplateNode{
		Token:    tok,
		Closed:   false,
		Ident:    ident,
		children: make([]Node, 0),
	}
}

func (tn *TemplateNode) Children() []Node {
	return tn.children
}

func (tn *TemplateNode) Append(node Node) {
	tn.children = append(tn.children, node)
}

func (tn *TemplateNode) String() string {
	return tn.Token.Literal
}

func (tn *TemplateNode) MarshalJSON() ([]byte, error) {
	n := struct {
		NodeType string
		Token    string
		Children []Node
		Closed   bool
	}{
		NodeType: "TemplateNode",
		Token:    tn.Token.Literal,
		Children: tn.children,
		Closed:   tn.Closed,
	}
	return json.Marshal(n)
}
