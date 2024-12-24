package main

import (
	"flag"
	"fmt"
	"iter"
	"slices"
	"strings"

	"github.com/bentekkie/advent_of_code_2024/pkg/bengraph"
	"github.com/bentekkie/advent_of_code_2024/pkg/inputs"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/topo"
	"gonum.org/v1/gonum/stat/combin"
)

func main() {
	flag.Parse()
	part1(inputs.Lines())
	part2(inputs.Lines())
}

func comb(n, k int) iter.Seq[[]int] {
	return func(yield func([]int) bool) {
		data := make([]int, k)
		for i := range data {
			data[i] = i
		}
		if !yield(data) {
			return
		}
		for range combin.Binomial(n, k) - 1 {
			next := make([]int, k)
			copy(next, data)
			for j := k - 1; j >= 0; j-- {
				if next[j] == n+j-k {
					continue
				}
				next[j]++
				for l := j + 1; l < k; l++ {
					next[l] = next[j] + l - j
				}
				break
			}
			if !yield(next) {
				return
			}
			data = next
		}
	}
}

func part1(input iter.Seq[string]) {
	g := simple.NewUndirectedGraph()
	compNames := map[string]*bengraph.Node[string]{}
	for line := range input {
		line := strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Split(line, "-")
		n1, ok := compNames[parts[0]]
		if !ok {
			n1 = bengraph.NewNode(g.NewNode(), parts[0])
			compNames[parts[0]] = n1
			g.AddNode(n1)
		}
		n2, ok := compNames[parts[1]]
		if !ok {
			n2 = bengraph.NewNode(g.NewNode(), parts[1])
			compNames[parts[1]] = n2
			g.AddNode(n2)
		}
		g.SetEdge(g.NewEdge(n1, n2))
	}
	sets := map[string]struct{}{}
	for _, ns := range topo.BronKerbosch(g) {
		if len(ns) < 3 {
			continue
		}
		for idxs := range comb(len(ns), 3) {
			ns := []string{
				ns[idxs[0]].(*bengraph.Node[string]).Data,
				ns[idxs[1]].(*bengraph.Node[string]).Data,
				ns[idxs[2]].(*bengraph.Node[string]).Data,
			}
			if ns[0][0] == 't' || ns[1][0] == 't' || ns[2][0] == 't' {
				slices.Sort(ns)
				sets[strings.Join(ns, "-")] = struct{}{}
			}
		}
	}
	fmt.Printf("Part 1: %v\n", len(sets))
}

func part2(input iter.Seq[string]) {
	g := simple.NewUndirectedGraph()
	compNames := map[string]*bengraph.Node[string]{}
	for line := range input {
		line := strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Split(line, "-")
		n1, ok := compNames[parts[0]]
		if !ok {
			n1 = bengraph.NewNode(g.NewNode(), parts[0])
			compNames[parts[0]] = n1
			g.AddNode(n1)
		}
		n2, ok := compNames[parts[1]]
		if !ok {
			n2 = bengraph.NewNode(g.NewNode(), parts[1])
			compNames[parts[1]] = n2
			g.AddNode(n2)
		}
		g.SetEdge(g.NewEdge(n1, n2))
	}
	var maxG []string
	for _, ns := range topo.BronKerbosch(g) {
		var group []string
		for _, n := range ns {
			group = append(group, n.(*bengraph.Node[string]).Data)
		}
		if len(group) > len(maxG) {
			maxG = group
		}
	}
	slices.Sort(maxG)
	fmt.Printf("Part 2: %v\n", strings.Join(maxG, ","))
}
