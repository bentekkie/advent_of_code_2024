package main

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/bentekkie/advent_of_code_2024/pkg/inputs"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/simple"
)

func main() {
	flag.Parse()
	part1(&inputs.Grid{})
	part2(&inputs.Grid{})
}

func part1(grid *inputs.Grid) {
	elevations := map[complex128]int{}
	nodes := map[complex128]graph.Node{}
	g := simple.NewDirectedGraph()
	starts := []complex128{}
	ends := []complex128{}
	for loc, c := range grid.Locs() {
		elevation, _ := strconv.Atoi(string(c))
		elevations[loc] = elevation
		nodes[loc] = g.NewNode()
		g.AddNode(nodes[loc])
		if elevation == 0 {
			starts = append(starts, loc)
		}
		if elevation == 9 {
			ends = append(ends, loc)
		}
	}
	for curr, elevation := range elevations {
		if el, ok := elevations[curr-1]; ok && el == elevation+1 {
			g.SetEdge(g.NewEdge(nodes[curr], nodes[curr-1]))
		}
		if el, ok := elevations[curr+1]; ok && el == elevation+1 {
			g.SetEdge(g.NewEdge(nodes[curr], nodes[curr+1]))
		}
		if el, ok := elevations[curr-1i]; ok && el == elevation+1 {
			g.SetEdge(g.NewEdge(nodes[curr], nodes[curr-1i]))
		}
		if el, ok := elevations[curr+1i]; ok && el == elevation+1 {
			g.SetEdge(g.NewEdge(nodes[curr], nodes[curr+1i]))
		}
	}
	allps := path.DijkstraAllPaths(g)
	s := 0
	for _, start := range starts {
		for _, end := range ends {
			ps, _ := allps.AllBetween(nodes[start].ID(), nodes[end].ID())
			if len(ps) > 0 {
				s += 1
			}
		}
	}
	fmt.Printf("Part 1: %d\n", s)
}

func part2(grid *inputs.Grid) {
	elevations := map[complex128]int{}
	nodes := map[complex128]graph.Node{}
	g := simple.NewDirectedGraph()
	starts := []complex128{}
	ends := []complex128{}
	for loc, c := range grid.Locs() {
		elevation, _ := strconv.Atoi(string(c))
		elevations[loc] = elevation
		nodes[loc] = g.NewNode()
		g.AddNode(nodes[loc])
		if elevation == 0 {
			starts = append(starts, loc)
		}
		if elevation == 9 {
			ends = append(ends, loc)
		}
	}
	for curr, elevation := range elevations {
		if el, ok := elevations[curr-1]; ok && el == elevation+1 {
			g.SetEdge(g.NewEdge(nodes[curr], nodes[curr-1]))
		}
		if el, ok := elevations[curr+1]; ok && el == elevation+1 {
			g.SetEdge(g.NewEdge(nodes[curr], nodes[curr+1]))
		}
		if el, ok := elevations[curr-1i]; ok && el == elevation+1 {
			g.SetEdge(g.NewEdge(nodes[curr], nodes[curr-1i]))
		}
		if el, ok := elevations[curr+1i]; ok && el == elevation+1 {
			g.SetEdge(g.NewEdge(nodes[curr], nodes[curr+1i]))
		}
	}
	allps := path.DijkstraAllPaths(g)
	s := 0
	for _, start := range starts {
		for _, end := range ends {
			ps, _ := allps.AllBetween(nodes[start].ID(), nodes[end].ID())
			s += len(ps)
		}
	}
	fmt.Printf("Part 2: %d\n", s)
}
