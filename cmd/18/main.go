package main

import (
	"flag"
	"fmt"
	"iter"
	"strconv"
	"strings"

	"github.com/bentekkie/advent_of_code_2024/pkg/bengraph"
	"github.com/bentekkie/advent_of_code_2024/pkg/inputs"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/topo"
)

func main() {
	flag.Parse()
	part1(inputs.Lines())
	part2(inputs.Lines())
}

func fullGrid(w, h int) iter.Seq[complex128] {
	return func(yield func(complex128) bool) {
		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				if !yield(complex(float64(x), float64(y))) {
					return
				}
			}
		}
	}
}

func mustParse[T ~int | ~float64](s string) T {
	n, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return T(n)
}

var dirs = [4]complex128{1, 1i, -1, -1i}

func part1(input iter.Seq[string]) {
	pts := map[complex128]struct{}{}
	for pt := range fullGrid(71, 71) {
		pts[pt] = struct{}{}
	}
	dropped := 0
	for line := range input {
		if dropped == 1024 {
			break
		}
		parts := strings.Split(strings.TrimSpace(line), ",")
		pt := complex(mustParse[float64](parts[0]), mustParse[float64](parts[1]))
		delete(pts, pt)
		dropped++
	}
	g := simple.NewUndirectedGraph()
	nodes := map[complex128]graph.Node{}
	for pt := range pts {
		n := g.NewNode()
		nodes[pt] = n
		g.AddNode(n)
	}
	for pt := range pts {
		for _, d := range dirs {
			if _, ok := nodes[pt+d]; ok {
				g.SetEdge(g.NewEdge(nodes[pt], nodes[pt+d]))
			}
		}
	}
	start := nodes[0]
	end := nodes[complex(70, 70)]
	shortest := path.DijkstraFrom(start, g)
	_, w := shortest.To(end.ID())
	fmt.Printf("Part 1: %v\n", w)
}

func part2(input iter.Seq[string]) {
	pts := map[complex128]struct{}{}
	for pt := range fullGrid(71, 71) {
		pts[pt] = struct{}{}
	}
	g := simple.NewUndirectedGraph()
	nodes := map[complex128]*bengraph.Node[complex128]{}
	for pt := range pts {
		n := bengraph.NewNode(g.NewNode(), pt)
		nodes[pt] = n
		g.AddNode(n)
	}
	for pt := range pts {
		for _, d := range dirs {
			if _, ok := nodes[pt+d]; ok {
				g.SetEdge(g.NewEdge(nodes[pt], nodes[pt+d]))
			}
		}
	}
	dropped := 0
	bad := ""
	start := nodes[0]
	end := nodes[complex(70, 70)]
	for line := range input {
		parts := strings.Split(strings.TrimSpace(line), ",")
		pt := complex(mustParse[float64](parts[0]), mustParse[float64](parts[1]))
		g.RemoveNode(nodes[pt].ID())
		dropped++
		if !topo.PathExistsIn(g, start, end) {
			bad = strings.TrimSpace(line)
			break
		}
	}
	fmt.Printf("Part 1: %s\n", bad)
}
