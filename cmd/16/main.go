package main

import (
	"flag"
	"fmt"

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

type State struct {
	loc complex128
	dir complex128
}

var dirs = []complex128{1, -1, 1i, -1i}

func part1(input *inputs.Grid) {
	var start, end complex128
	openStates := map[State]graph.Node{}
	g := simple.NewWeightedUndirectedGraph(0, 0)
	for loc, c := range input.Locs() {
		if c == 'S' {
			start = loc
		} else if c == 'E' {
			end = loc
		}
		if c == '.' || c == 'S' || c == 'E' {
			for _, dir := range dirs {
				openStates[State{loc, dir}] = g.NewNode()
				g.AddNode(openStates[State{loc, dir}])
			}
			for i := range dirs {
				g.SetWeightedEdge(g.NewWeightedEdge(openStates[State{loc, dirs[i]}], openStates[State{loc, dirs[(i+1)%4]}], 1000))
			}
		}
	}
	startNode := openStates[State{loc: start, dir: 1}]
	endNode := g.NewNode()
	g.AddNode(endNode)
	for _, dir := range dirs {
		g.SetWeightedEdge(g.NewWeightedEdge(openStates[State{loc: end, dir: dir}], endNode, 0))
	}
	for state, currNode := range openStates {
		forward := State{
			dir: state.dir,
			loc: state.loc + state.dir,
		}
		if forwardNode, ok := openStates[forward]; ok {
			g.SetWeightedEdge(g.NewWeightedEdge(currNode, forwardNode, 1))
		}
	}
	p, _ := path.AStar(startNode, endNode, g, nil)
	p.WeightTo(endNode.ID())
	// TODO: implement me
	fmt.Printf("Part 1: %v\n", p.WeightTo(endNode.ID()))
}

func part2(input *inputs.Grid) {
	var start, end complex128
	openStates := map[State]graph.Node{}
	nodeToLoc := map[int64]complex128{}
	g := simple.NewWeightedUndirectedGraph(0, 0)
	for loc, c := range input.Locs() {
		if c == 'S' {
			start = loc
		} else if c == 'E' {
			end = loc
		}
		if c == '.' || c == 'S' || c == 'E' {
			for _, dir := range dirs {
				openStates[State{loc, dir}] = g.NewNode()
				g.AddNode(openStates[State{loc, dir}])
				nodeToLoc[openStates[State{loc, dir}].ID()] = loc
			}
			for i := range dirs {
				g.SetWeightedEdge(g.NewWeightedEdge(openStates[State{loc, dirs[i]}], openStates[State{loc, dirs[(i+1)%4]}], 1000))
			}
		}
	}
	startNode := openStates[State{loc: start, dir: 1}]
	endNode := g.NewNode()
	g.AddNode(endNode)
	nodeToLoc[endNode.ID()] = end
	for _, dir := range dirs {
		g.SetWeightedEdge(g.NewWeightedEdge(openStates[State{loc: end, dir: dir}], endNode, 0))
	}
	for state, currNode := range openStates {
		forward := State{
			dir: state.dir,
			loc: state.loc + state.dir,
		}
		if forwardNode, ok := openStates[forward]; ok {
			g.SetWeightedEdge(g.NewWeightedEdge(currNode, forwardNode, 1))
		}
	}
	all, _ := path.BellmanFordAllFrom(startNode, g)
	ps, _ := all.AllTo(endNode.ID())
	pathLocs := map[complex128]struct{}{}
	for _, p := range ps {
		for _, n := range p {
			pathLocs[nodeToLoc[n.ID()]] = struct{}{}
		}
	}

	fmt.Printf("Part 2: %v\n", len(pathLocs))
}
