package main

import (
	"flag"
	"fmt"
	"iter"
	"math"

	"github.com/bentekkie/advent_of_code_2024/pkg/bengraph"
	"github.com/bentekkie/advent_of_code_2024/pkg/flags"
	"github.com/bentekkie/advent_of_code_2024/pkg/inputs"
	"gonum.org/v1/gonum/graph/path"
)

func main() {
	flag.Parse()
	defer flags.CPUProfile()()
	part1(&inputs.Grid{})
	part2(&inputs.Grid{})
}

var dirs = [4]complex128{1, -1, 1i, -1i}

func part1(input *inputs.Grid) {
	var start, end complex128
	walls := map[complex128]struct{}{}
	for loc, c := range input.Locs() {
		if c == 'S' {
			start = loc
		}
		if c == 'E' {
			end = loc
		}
		if c == '#' {
			walls[loc] = struct{}{}
		}
	}
	g := &bengraph.Grid{
		Missing: walls,
		Max:     input.Max(),
	}
	allFromStart := path.DijkstraFrom(g.Node(start), g)
	allFromEnd := path.DijkstraFrom(g.Node(end), g)
	p, w := allFromStart.To(g.LocToID(end))
	s := 0
	for _, n := range bengraph.Path[complex128](p) {
		for _, d := range dirs {
			wall := n.Data + d
			if _, ok := walls[wall]; ok {
				for _, d := range dirs {
					next := wall + d
					if nextNode := g.Node(next); nextNode != nil && next != n.Data {
						dist := allFromStart.WeightTo(n.ID()) + allFromEnd.WeightTo(nextNode.ID()) + 2
						if dist <= w-100 {
							s++
						}
					}
				}
			}
		}
	}
	fmt.Printf("Part 1: %v\n", s)
}

func part2(input *inputs.Grid) {
	var start, end complex128
	walls := map[complex128]struct{}{}
	for loc, c := range input.Locs() {
		if c == 'S' {
			start = loc
		}
		if c == 'E' {
			end = loc
		}
		if c == '#' {
			walls[loc] = struct{}{}
		}
	}
	g := &bengraph.Grid{
		Missing: walls,
		Max:     input.Max(),
	}
	allFromStart := path.DijkstraFrom(g.Node(start), g)
	allFromEnd := path.DijkstraFrom(g.Node(end), g)
	p, w := allFromStart.To(g.LocToID(end))
	s := 0
	for _, n := range bengraph.Path[complex128](p) {
		for next := range around(n.Data, 21) {
			if nextNode := g.Node(next); nextNode != nil && next != n.Data {
				dist := allFromStart.WeightTo(n.ID()) + allFromEnd.WeightTo(nextNode.ID()) + cdist(n.Data, next)
				if dist <= w-100 {
					s++
				}
			}
		}
	}
	fmt.Printf("Part 2: %v\n", s)
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func cdist(a, b complex128) float64 {
	return math.Abs(real(a)-real(b)) + math.Abs(imag(a)-imag(b))
}

func around(loc complex128, rad int) iter.Seq[complex128] {
	return func(yield func(complex128) bool) {
		for i := -rad + 1; i < rad; i++ {
			maxj := rad - abs(i)
			for j := -maxj + 1; j < maxj; j++ {
				if !yield(loc + complex(float64(i), float64(j))) {
					return
				}
			}
		}
	}
}
