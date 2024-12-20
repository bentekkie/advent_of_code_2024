package bengraph

import (
	"iter"

	"gonum.org/v1/gonum/graph"
)

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

func Nodes(ns graph.Nodes) iter.Seq[graph.Node] {
	return func(yield func(graph.Node) bool) {
		for ns.Next() {
			if !yield(ns.Node()) {
				return
			}
		}
	}
}

func DataToNodes[T comparable](g graph.Graph) map[T]*Node[T] {
	nodes := map[T]*Node[T]{}
	for n := range Nodes(g.Nodes()) {
		nodes[n.(*Node[T]).Data] = n.(*Node[T])
	}
	return nodes
}

func Path[T comparable](p []graph.Node) []*Node[T] {
	data := make([]*Node[T], len(p))
	for i, n := range p {
		data[i] = n.(*Node[T])
	}
	return data
}
