package bengraph

import "gonum.org/v1/gonum/graph"

type Node[T any] struct {
	graph.Node
	Data T
}

func NewNode[T any](n graph.Node, data T) *Node[T] {
	return &Node[T]{
		Node: n,
		Data: data,
	}
}
