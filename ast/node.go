package ast

// Node represents a node in an abstract syntax tree.
type Node interface {
	// String returns a string representation of the node.
	String() string
	// Children returns the children nodes of the node.
	Children() []Node
	// Append appends a child node to the node.
	Append(node Node)
}
