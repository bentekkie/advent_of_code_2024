package main

import (
	"flag"
	"fmt"
	"iter"

	"github.com/bentekkie/advent_of_code_2024/pkg/inputs"
)

func main() {
	flag.Parse()
	part1(inputs.Grid())
	part2(inputs.Grid())
}

func part1(grid iter.Seq2[complex128, string]) {
	annennas := map[string][]complex128{}
	boundX, boundY := 0, 0
	for loc, c := range grid {
		if c != "." {
			annennas[c] = append(annennas[c], loc)
		}
		if real(loc) > float64(boundX) {
			boundX = int(real(loc))
		}
		if imag(loc) > float64(boundY) {
			boundY = int(imag(loc))
		}
	}
	ps := map[complex128]struct{}{}
	for _, locs := range annennas {
		for _, locA := range locs {
			for _, locB := range locs {
				if locA != locB {
					diff := locA - locB
					p1 := locB - diff
					p2 := locA + diff
					if real(p1) >= 0 && real(p1) <= float64(boundX) && imag(p1) >= 0 && imag(p1) <= float64(boundY) {
						ps[p1] = struct{}{}
					}
					if real(p2) >= 0 && real(p2) <= float64(boundX) && imag(p2) >= 0 && imag(p2) <= float64(boundY) {
						ps[p2] = struct{}{}
					}
				}
			}
		}
	}
	fmt.Printf("Part 1: %d\n", len(ps))
}

func part2(grid iter.Seq2[complex128, string]) {
	annennas := map[string][]complex128{}
	boundX, boundY := 0, 0
	for loc, c := range grid {
		if c != "." {
			annennas[c] = append(annennas[c], loc)
		}
		if real(loc) > float64(boundX) {
			boundX = int(real(loc))
		}
		if imag(loc) > float64(boundY) {
			boundY = int(imag(loc))
		}
	}
	inBounds := func(c complex128) bool {
		return real(c) >= 0 && real(c) <= float64(boundX) && imag(c) >= 0 && imag(c) <= float64(boundY)
	}
	ps := map[complex128]struct{}{}
	for _, locs := range annennas {
		for _, locA := range locs {
			for _, locB := range locs {
				if locA != locB {
					diff := locA - locB
					g := GCD(real(diff), imag(diff))
					diff /= complex(g, 0)
					p := locB
					for inBounds(p) {
						ps[p] = struct{}{}
						p += diff
					}
					p = locB
					for inBounds(p) {
						ps[p] = struct{}{}
						p -= diff
					}
				}
			}
		}
	}
	fmt.Printf("Part 1: %d\n", len(ps))
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func GCD(af, bf float64) float64 {
	a, b := abs(int(af)), abs(int(bf))
	if a == 0 {
		return bf
	}
	if b == 0 {
		return af
	}
	for a != b {
		if a > b {
			a -= b
		} else {
			b -= a
		}
	}

	return float64(a)
}
