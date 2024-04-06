package ast

import (
	"encoding/json"

	"github.com/connormckelvey/tmplrun/token"
)

// TextNode represents a node in the abstract syntax tree corresponding to text.
type TextNode struct {
	// Token is the token associated with the text node.
	Token *token.Token
}

// String returns the string representation of the text node.
func (tn *TextNode) String() string {
	return tn.Token.Literal
}

// Children returns nil for text nodes as they do not have children.
func (tn *TextNode) Children() []Node {
	return nil
}

// Append panics when attempting to append a child node to a text node, as text nodes cannot have children.
func (tn *TextNode) Append(node Node) {
	panic("cannot append child to text node")
}

// MarshalJSON returns the JSON encoding of the text node.
func (tn *TextNode) MarshalJSON() ([]byte, error) {
	dd := struct {
		NodeType string
		Text     string
	}{
		NodeType: "TextNode",
		Text:     tn.Token.Literal,
	}
	return json.Marshal(dd)
}
