package main

import (
	"flag"
	"fmt"

	"github.com/bentekkie/advent_of_code_2024/pkg/benlog"
	"github.com/bentekkie/advent_of_code_2024/pkg/inputs"
)

func main() {
	flag.Parse()
	benlog.Timed(func() {
		part1(&inputs.Grid{})
	})
	benlog.Timed(func() {
		part2(&inputs.Grid{})
	})
}

func part1(grid *inputs.Grid) {
	annennas := map[rune][]complex128{}
	var cnt int
	for loc, c := range grid.Locs() {
		if c != '.' {
			annennas[c] = append(annennas[c], loc)
		}
		cnt++
	}
	bound := grid.Max()
	ps := make(map[complex128]struct{}, cnt)
	for _, locs := range annennas {
		for i, locA := range locs {
			for _, locB := range locs[i+1:] {
				if p := 2*locB - locA; real(p) >= 0 && real(p) <= real(bound) && imag(p) >= 0 && imag(p) <= imag(bound) {
					ps[p] = struct{}{}
				}
				if p := 2*locA - locB; real(p) >= 0 && real(p) <= real(bound) && imag(p) >= 0 && imag(p) <= imag(bound) {
					ps[p] = struct{}{}
				}
			}
		}
	}
	fmt.Printf("Part 1: %d\n", len(ps))
}

func part2(grid *inputs.Grid) {
	annennas := map[rune][]complex128{}
	var cnt int
	for loc, c := range grid.Locs() {
		if c != '.' {
			annennas[c] = append(annennas[c], loc)
		}
		cnt++
	}
	bound := grid.Max()
	ps := make(map[complex128]struct{}, cnt)
	for _, locs := range annennas {
		for i, locA := range locs {
			for _, locB := range locs[i+1:] {
				diff := locA - locB
				for p := locB; real(p) >= 0 && real(p) <= real(bound) && imag(p) >= 0 && imag(p) <= imag(bound); p += diff {
					ps[p] = struct{}{}
				}
				for p := locB - diff; real(p) >= 0 && real(p) <= real(bound) && imag(p) >= 0 && imag(p) <= imag(bound); p -= diff {
					ps[p] = struct{}{}
				}
			}
		}
	}
	fmt.Printf("Part 2: %d\n", len(ps))
}
