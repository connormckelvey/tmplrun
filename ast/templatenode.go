package ast

import (
	"encoding/json"

	"github.com/connormckelvey/tmplrun/lexer"
	"github.com/connormckelvey/tmplrun/token"
)

// TemplateNode represents a node in the abstract syntax tree corresponding to a template tag.
type TemplateNode struct {
	// Token is the token associated with the template node.
	Token *token.Token
	// Closed indicates whether the template tag is closed.
	Closed bool
	// Ident is the identifier of the template tag.
	Ident string
	// children is the slice of child nodes.
	children []Node
}

func parseTemplateTag(lit string) (ok bool, ident string, padding string) {
	match := lexer.OpenTagPattern.FindStringSubmatch(lit)
	if match == nil {
		return false, "", ""
	}
	return true, match[2], match[3]
}

// NewTemplateNode creates a new TemplateNode instance from the given token.
func NewTemplateNode(tok *token.Token) *TemplateNode {
	_, ident, _ := parseTemplateTag(tok.Literal)
	return &TemplateNode{
		Token:    tok,
		Closed:   false,
		Ident:    ident,
		children: make([]Node, 0),
	}
}

// Children returns the children nodes of the template node.
func (tn *TemplateNode) Children() []Node {
	return tn.children
}

// Append appends a child node to the template node.
func (tn *TemplateNode) Append(node Node) {
	tn.children = append(tn.children, node)
}

// String returns the string representation of the template node.
func (tn *TemplateNode) String() string {
	return tn.Token.Literal
}

// MarshalJSON returns the JSON encoding of the template node.
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
