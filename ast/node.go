package ast

type Node interface {
	String() string
	Children() []Node
	Append(node Node)
}
