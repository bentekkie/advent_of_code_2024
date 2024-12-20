package bengraph

import (
	"iter"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/iterator"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/traverse"
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

type Grid struct {
	Missing map[complex128]struct{}
	Max     complex128
}

var _ traverse.Graph = (*Grid)(nil)

var dirs = [4]complex128{1, -1, 1i, -1i}

func (g *Grid) IdToLoc(id int64) complex128 {
	return complex(float64(id%int64(real(g.Max)+1)), float64(id/int64(real(g.Max)+1)))
}

func (g *Grid) LocToID(loc complex128) int64 {
	return int64(real(loc) + imag(loc)*(real(g.Max)+1))
}

func (g *Grid) Node(loc complex128) graph.Node {
	if !g.valid(loc) {
		return nil
	}
	return NewNode(simple.Node(g.LocToID(loc)), loc)
}

func (g *Grid) valid(loc complex128) bool {
	if real(loc) >= 0 && real(loc) <= real(g.Max) && imag(loc) >= 0 && imag(loc) <= imag(g.Max) {
		if _, ok := g.Missing[loc]; !ok {
			return true
		}
	}
	return false
}

func (g *Grid) From(id int64) graph.Nodes {
	loc := g.IdToLoc(id)
	var nodes []graph.Node
	for _, d := range dirs {
		next := loc + d
		if g.valid(next) {
			nodes = append(nodes, g.Node(next))
		}
	}
	if nodes == nil {
		return graph.Empty
	}
	return iterator.NewOrderedNodes(nodes)
}

func (g *Grid) Edge(uid, vid int64) graph.Edge {
	u, v := g.IdToLoc(uid), g.IdToLoc(vid)
	if !g.valid(u) || !g.valid(v) {
		return nil
	}
	if u+1 == v || u-1 == v || u+1i == v || u-1i == v {
		return simple.Edge{F: g.Node(u), T: g.Node(v)}
	}
	return nil
}
