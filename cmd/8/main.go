package main

import (
	"flag"
	"fmt"
	"iter"
	"math"

	"github.com/bentekkie/advent_of_code_2024/pkg/inputs"
)

func main() {
	flag.Parse()
	part1(inputs.Grid())
	part2(inputs.Grid())
}

func part1(grid iter.Seq2[complex128, string]) {
	annennas := map[string][]complex128{}
	var boundX, boundY float64
	for loc, c := range grid {
		if c != "." {
			annennas[c] = append(annennas[c], loc)
		}
		if real(loc) > boundX {
			boundX = real(loc)
		}
		if imag(loc) > boundY {
			boundY = imag(loc)
		}
	}
	inBounds := func(c complex128) bool {
		return real(c) >= 0 && real(c) <= boundX && imag(c) >= 0 && imag(c) <= boundY
	}
	ps := map[complex128]struct{}{}
	for _, locs := range annennas {
		for _, locA := range locs {
			for _, locB := range locs {
				if locA != locB {
					if p := 2*locB - locA; inBounds(p) {
						ps[p] = struct{}{}
					}
					if p := 2*locA - locB; inBounds(p) {
						ps[p] = struct{}{}
					}
				}
			}
		}
	}
	fmt.Printf("Part 1: %d\n", len(ps))
}

func part2(grid iter.Seq2[complex128, string]) {
	annennas := map[string][]complex128{}
	var boundX, boundY float64
	for loc, c := range grid {
		if c != "." {
			annennas[c] = append(annennas[c], loc)
		}
		if real(loc) > boundX {
			boundX = real(loc)
		}
		if imag(loc) > boundY {
			boundY = imag(loc)
		}
	}
	inBounds := func(c complex128) bool {
		return real(c) >= 0 && real(c) <= boundX && imag(c) >= 0 && imag(c) <= boundY
	}
	ps := map[complex128]struct{}{}
	for _, locs := range annennas {
		for _, locA := range locs {
			for _, locB := range locs {
				if locA != locB {
					diff := locA - locB
					for p := locB; inBounds(p); p += diff {
						ps[p] = struct{}{}
					}
					for p := locB; inBounds(p); p -= diff {
						ps[p] = struct{}{}
					}
				}
			}
		}
	}
	fmt.Printf("Part 1: %d\n", len(ps))
}

func GCD(a, b float64) float64 {
	a, b = math.Abs(a), math.Abs(b)
	if a == 0 {
		return b
	}
	if b == 0 {
		return a
	}
	for a != b {
		if a > b {
			a -= b
		} else {
			b -= a
		}
	}
	return a
}
