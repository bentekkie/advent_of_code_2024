package main

import (
	"flag"
	"fmt"
	"iter"
	"sort"
	"strings"

	"github.com/bentekkie/advent_of_code_2024/pkg/bengraph"
	"github.com/bentekkie/advent_of_code_2024/pkg/benlog"
	"github.com/bentekkie/advent_of_code_2024/pkg/inputs"
	"github.com/bentekkie/advent_of_code_2024/pkg/parse"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/simple"
)

func main() {
	flag.Parse()
	part1(inputs.Lines())
	benlog.Timed(func() {
		part2try2(inputs.Lines())
	})
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

var dirs = [4]complex128{1, 1i, -1, -1i}

func part1(input iter.Seq[string]) {
	dropped := 0
	missing := map[complex128]struct{}{}
	for line := range input {
		if dropped == 1024 {
			break
		}
		parts := parse.NumList[float64](strings.TrimSpace(line), ",")
		pt := complex(parts[0], parts[1])
		missing[pt] = struct{}{}
		dropped++
	}
	v, i := bfs(complex(0, 0), complex(70, 70), missing)
	fmt.Printf("Part 1: %t %v\n", v, i)
}

func makeGraph(toRemove []complex128) (graph.Node, graph.Node, graph.Graph) {
	pts := map[complex128]struct{}{}
	for pt := range fullGrid(71, 71) {
		pts[pt] = struct{}{}
	}
	for _, pt := range toRemove {
		delete(pts, pt)
	}
	g := simple.NewUndirectedGraph()
	for pt := range pts {
		g.AddNode(bengraph.NewNode(g.NewNode(), pt))
	}
	nodes := bengraph.DataToNodes[complex128](g)
	for pt := range pts {
		for _, d := range dirs {
			if _, ok := nodes[pt+d]; ok {
				g.SetEdge(g.NewEdge(nodes[pt], nodes[pt+d]))
			}
		}
	}
	start := nodes[0]
	end := nodes[complex(70, 70)]
	return start, end, g
}

func part2(input iter.Seq[string]) {
	pts := map[complex128]struct{}{}
	for pt := range fullGrid(71, 71) {
		pts[pt] = struct{}{}
	}
	g := simple.NewUndirectedGraph()
	for pt := range pts {
		g.AddNode(bengraph.NewNode(g.NewNode(), pt))
	}
	nodes := bengraph.DataToNodes[complex128](g)
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
	pm := make(map[complex128]struct{}, 70*70)
	nextShortest := func() bool {
		p, _ := path.AStar(start, end, g, nil)
		ns, _ := p.To(end.ID())
		if ns == nil {
			return false
		}
		clear(pm)
		for _, n := range ns {
			pm[n.(*bengraph.Node[complex128]).Data] = struct{}{}
		}
		return true
	}
	nextShortest()
	for line := range input {
		parts := parse.NumList[float64](strings.TrimSpace(line), ",")
		pt := complex(parts[0], parts[1])
		g.RemoveNode(nodes[pt].ID())
		if _, ok := pm[pt]; ok && dropped >= 70 {
			if !nextShortest() {
				bad = strings.TrimSpace(line)
				break
			}
		}
		dropped++
	}
	simple.NewDirectedGraph()
	fmt.Printf("Part 2: %s\n", bad)
}

func bfs(start, target complex128, missing map[complex128]struct{}) (bool, int) {
	toVisit := map[complex128]struct{}{start: {}}
	visited := map[complex128]struct{}{}
	steps := 0
	for len(toVisit) > 0 {
		nextToVisit := map[complex128]struct{}{}
		for curr := range toVisit {
			if curr == target {
				return true, steps
			}
			visited[curr] = struct{}{}
			for _, d := range dirs {
				next := curr + d
				if _, ok := missing[next]; ok {
					continue
				}
				if _, ok := visited[next]; ok {
					continue
				}
				if _, ok := toVisit[next]; ok {
					continue
				}
				if real(next) < 0 || real(next) > 70 || imag(next) < 0 || imag(next) > 70 {
					continue
				}
				nextToVisit[next] = struct{}{}
			}
		}
		steps++
		toVisit = nextToVisit
	}
	return false, 0
}

type ginfo struct {
	g          graph.Graph
	start, end graph.Node
}

func part2try2(input iter.Seq[string]) {
	missing := map[complex128]int{}
	atTime := map[int]string{}
	time := 0
	for line := range input {
		parts := parse.NumList[float64](strings.TrimSpace(line), ",")
		pt := complex(parts[0], parts[1])
		missing[pt] = time
		atTime[time] = strings.TrimSpace(line)
		time++
	}
	k := sort.Search(len(missing), func(i int) bool {
		v, _ := bfsuntil(complex(0, 0), complex(70, 70), missing, i)
		return !v
	})
	fmt.Printf("Part 2: %s\n", atTime[k])
}

func bfsuntil(start, target complex128, missing map[complex128]int, until int) (bool, int) {
	toVisit := map[complex128]struct{}{start: {}}
	visited := map[complex128]struct{}{}
	steps := 0
	for len(toVisit) > 0 {
		nextToVisit := map[complex128]struct{}{}
		for curr := range toVisit {
			if curr == target {
				return true, steps
			}
			visited[curr] = struct{}{}
			for _, d := range dirs {
				next := curr + d
				if time, ok := missing[next]; ok && time <= until {
					continue
				}
				if _, ok := visited[next]; ok {
					continue
				}
				if _, ok := toVisit[next]; ok {
					continue
				}
				if real(next) < 0 || real(next) > 70 || imag(next) < 0 || imag(next) > 70 {
					continue
				}
				nextToVisit[next] = struct{}{}
			}
		}
		steps++
		toVisit = nextToVisit
	}
	return false, 0
}
