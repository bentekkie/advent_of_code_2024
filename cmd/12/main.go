package main

import (
	"flag"
	"fmt"

	"github.com/bentekkie/advent_of_code_2024/pkg/inputs"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/topo"
)

func main() {
	flag.Parse()
	part1(&inputs.Grid{})
	part2(&inputs.Grid{})
}

func part1(input *inputs.Grid) {
	locs := map[complex128]rune{}
	g := simple.NewUndirectedGraph()
	nodes := map[complex128]graph.Node{}
	locsFromNodeIds := map[int64]complex128{}
	for loc, c := range input.Locs() {
		locs[loc] = c
		nodes[loc] = g.NewNode()
		g.AddNode(nodes[loc])
		locsFromNodeIds[nodes[loc].ID()] = loc
	}
	for loc, c := range locs {
		if other, ok := locs[loc-1]; ok && c == other {
			g.SetEdge(g.NewEdge(nodes[loc], nodes[loc-1]))
		}
		if other, ok := locs[loc+1]; ok && c == other {
			g.SetEdge(g.NewEdge(nodes[loc], nodes[loc+1]))
		}
		if other, ok := locs[loc-1i]; ok && c == other {
			g.SetEdge(g.NewEdge(nodes[loc], nodes[loc-1i]))
		}
		if other, ok := locs[loc+1i]; ok && c == other {
			g.SetEdge(g.NewEdge(nodes[loc], nodes[loc+1i]))
		}
	}
	s := 0
	for _, component := range topo.ConnectedComponents(g) {
		area := len(component)
		compLocs := make(map[complex128]struct{}, len(component))
		for _, n := range component {
			compLocs[locsFromNodeIds[n.ID()]] = struct{}{}
		}
		perim := 0
		for loc := range compLocs {
			if _, ok := compLocs[loc-1]; !ok {
				perim++
			}
			if _, ok := compLocs[loc+1]; !ok {
				perim++
			}
			if _, ok := compLocs[loc-1i]; !ok {
				perim++
			}
			if _, ok := compLocs[loc+1i]; !ok {
				perim++
			}
		}
		s += area * perim
	}
	// TODO: implement me
	fmt.Printf("Part 1: %d\n", s)
}

type Dir int

const (
	V Dir = iota
	H
)

type side struct {
	a, b   complex128
	inside complex128
	dir    Dir
}

func newSide(a, b, inside complex128) side {
	if real(a) == real(b) {
		if imag(a) < imag(b) {
			return side{a, b, inside, V}
		} else {
			return side{b, a, inside, V}
		}
	} else {
		if real(a) < real(b) {
			return side{a, b, inside, H}
		} else {
			return side{b, a, inside, H}
		}
	}
}

func (s side) nextTo(o side) bool {
	if s.dir != o.dir {
		return false
	}
	if s.dir == H {
		return (s.a == o.a+1i && s.b == o.b+1i && s.inside == o.inside+1i) || (s.a == o.a-1i && s.b == o.b-1i && s.inside == o.inside-1i)
	} else {
		return (s.a == o.a+1 && s.b == o.b+1 && s.inside == o.inside+1) || (s.a == o.a-1 && s.b == o.b-1 && s.inside == o.inside-1)
	}
}

func part2(input *inputs.Grid) {
	locs := map[complex128]rune{}
	g := simple.NewUndirectedGraph()
	nodes := map[complex128]graph.Node{}
	locsFromNodeIds := map[int64]complex128{}
	for loc, c := range input.Locs() {
		locs[loc] = c
		nodes[loc] = g.NewNode()
		g.AddNode(nodes[loc])
		locsFromNodeIds[nodes[loc].ID()] = loc
	}
	for loc, c := range locs {
		if other, ok := locs[loc-1]; ok && c == other {
			g.SetEdge(g.NewEdge(nodes[loc], nodes[loc-1]))
		}
		if other, ok := locs[loc+1]; ok && c == other {
			g.SetEdge(g.NewEdge(nodes[loc], nodes[loc+1]))
		}
		if other, ok := locs[loc-1i]; ok && c == other {
			g.SetEdge(g.NewEdge(nodes[loc], nodes[loc-1i]))
		}
		if other, ok := locs[loc+1i]; ok && c == other {
			g.SetEdge(g.NewEdge(nodes[loc], nodes[loc+1i]))
		}
	}
	s := 0
	for _, component := range topo.ConnectedComponents(g) {
		area := len(component)
		compLocs := make(map[complex128]struct{}, len(component))
		for _, n := range component {
			compLocs[locsFromNodeIds[n.ID()]] = struct{}{}
		}
		sides := map[side]struct{}{}
		for loc := range compLocs {
			if _, ok := compLocs[loc-1]; !ok {
				sides[newSide(loc, loc-1, loc)] = struct{}{}
			}
			if _, ok := compLocs[loc+1]; !ok {
				sides[newSide(loc, loc+1, loc)] = struct{}{}
			}
			if _, ok := compLocs[loc-1i]; !ok {
				sides[newSide(loc, loc-1i, loc)] = struct{}{}
			}
			if _, ok := compLocs[loc+1i]; !ok {
				sides[newSide(loc, loc+1i, loc)] = struct{}{}
			}
		}
		sg := simple.NewUndirectedGraph()
		snodes := make(map[side]graph.Node, len(sides))
		for s := range sides {
			n := sg.NewNode()
			sg.AddNode(n)
			snodes[s] = n
		}
		for s := range sides {
			for o := range sides {
				if s != o && s.nextTo(o) {
					sg.SetEdge(sg.NewEdge(snodes[s], snodes[o]))
				}
			}
		}
		numSides := len(topo.ConnectedComponents(sg))

		s += area * numSides
	}
	fmt.Printf("Part 2: %d\n", s)
}
