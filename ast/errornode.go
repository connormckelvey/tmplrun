package ast

import (
	"encoding/json"

	"github.com/connormckelvey/tmplrun/token"
)

type ErrorNode struct {
	Token *token.Token
	Err   error
}

func (tn *ErrorNode) String() string {
	return tn.Token.Literal
}
func (tn *ErrorNode) Children() []Node {
	return nil
}
func (tn *ErrorNode) Append(node Node) {
	panic("cannot append child to text node")
}

func (tn *ErrorNode) MarshalJSON() ([]byte, error) {
	dd := struct {
		NodeType string
		Text     string
		Error    string
	}{
		NodeType: "ErrorNode",
		Text:     tn.Token.Literal,
		Error:    tn.Err.Error(),
	}
	return json.Marshal(dd)
}
