package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/bentekkie/advent_of_code_2024/pkg/inputs"
)

func main() {
	flag.Parse()
	input := inputs.String()
	part1(input)
	part2(input)
}

type KeyPins struct {
	p1, p2, p3, p4, p5 uint8
}

func (k KeyPins) overlaps(o KeyPins) bool {
	if k.p1+o.p1 > 5 {
		return false
	}
	if k.p2+o.p2 > 5 {
		return false
	}
	if k.p3+o.p3 > 5 {
		return false
	}
	if k.p4+o.p4 > 5 {
		return false
	}
	if k.p5+o.p5 > 5 {
		return false
	}
	return true
}

func (k KeyPins) String() string {
	return fmt.Sprintf("%d%d%d%d%d", k.p1, k.p2, k.p3, k.p4, k.p5)
}

func part1(input string) {
	locks := map[KeyPins]int{}
	keys := map[KeyPins]int{}
	for _, grid := range strings.Split(strings.TrimSpace(input), "\n\n") {
		if strings.HasPrefix(grid, "#####") {
			locks[parsePins(grid)]++
		} else {
			keys[parsePins(grid)]++
		}
	}
	s := 0
	for l, li := range locks {
		for k, ki := range keys {
			if l.overlaps(k) {
				s += li * ki
			}
		}
	}
	// TODO: implement me
	fmt.Printf("Part 1: %v\n", s)
}

func parsePins(grid string) (k KeyPins) {
	for _, line := range strings.Split(grid, "\n")[1:6] {
		line = strings.TrimSpace(line)
		if line[0] == '#' {
			k.p1++
		}
		if line[1] == '#' {
			k.p2++
		}
		if line[2] == '#' {
			k.p3++
		}
		if line[3] == '#' {
			k.p4++
		}
		if line[4] == '#' {
			k.p5++
		}
	}
	return
}

func part2(input string) {
	// TODO: implement me
	fmt.Printf("Part 2: FREE SQUARE!\n")
}
