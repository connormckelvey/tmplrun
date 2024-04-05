package ast

import (
	"encoding/json"
	"fmt"
)

type Document struct {
	children []Node
}

func (d *Document) Children() []Node {
	return d.children
}

func (d *Document) Append(node Node) {
	d.children = append(d.children, node)
}

func (d *Document) String() string {
	return fmt.Sprint(d.Children())
}

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
