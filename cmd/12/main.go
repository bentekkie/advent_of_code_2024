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
	in, out complex128
}

func (s side) dir() Dir {
	if real(s.in) == real(s.out) {
		return V
	}
	return H
}

func (s side) nextTo(o side) bool {
	if s.dir() != o.dir() {
		return false
	}
	if s.dir() == H {
		return (s.in == o.in+1i && s.out == o.out+1i) || (s.in == o.in-1i && s.out == o.out-1i)
	} else {
		return (s.in == o.in+1 && s.out == o.out+1) || (s.in == o.in-1 && s.out == o.out-1)
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
				sides[side{in: loc, out: loc - 1}] = struct{}{}
			}
			if _, ok := compLocs[loc+1]; !ok {
				sides[side{in: loc, out: loc + 1}] = struct{}{}
			}
			if _, ok := compLocs[loc-1i]; !ok {
				sides[side{in: loc, out: loc - 1i}] = struct{}{}
			}
			if _, ok := compLocs[loc+1i]; !ok {
				sides[side{in: loc, out: loc + 1i}] = struct{}{}
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
