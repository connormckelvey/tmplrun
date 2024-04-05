package ast

import (
	"encoding/json"

	"github.com/connormckelvey/tmplrun/token"
)

type TextNode struct {
	Token *token.Token
}

func (tn *TextNode) String() string {
	return tn.Token.Literal
}
func (tn *TextNode) Children() []Node {
	return nil
}
func (tn *TextNode) Append(node Node) {
	panic("cannot append child to text node")
}

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
