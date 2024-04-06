package ast

import (
	"encoding/json"
	"fmt"
)

// Document represents an abstract syntax tree document.
type Document struct {
	children []Node // children is the slice of nodes in the document.
}

// Children returns the children nodes of the document.
func (d *Document) Children() []Node {
	return d.children
}

// Append appends a node to the list of children nodes in the document.
func (d *Document) Append(node Node) {
	d.children = append(d.children, node)
}

// String returns a string representation of the document.
func (d *Document) String() string {
	return fmt.Sprint(d.Children())
}

// MarshalJSON returns the JSON encoding of the document.
func (d *Document) MarshalJSON() ([]byte, error) {
	dd := struct {
		NodeType string
		Children []Node
	}{
		NodeType: "Document",
		Children: d.children,
	}
	return json.Marshal(dd)
}
