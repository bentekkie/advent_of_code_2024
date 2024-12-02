package main

import (
	"flag"
	"fmt"
	"iter"
	"strconv"
	"strings"

	"github.com/bentekkie/advent_of_code_2024/pkg/inputs"
)

func main() {
	flag.Parse()
	part1(inputs.Lines())
	part2(inputs.Lines())
}

func part1(input iter.Seq[string]) {
	s := 0
	for line := range input {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if isSafe(parseLine(line)) {
			s++
		}
	}
	// TODO: implement me
	fmt.Printf("Part 1: %d\n", s)
}

func parseLine(line string) []int {
	parts := strings.Split(strings.TrimSpace(line), " ")
	levels := make([]int, len(parts))
	for i, part := range parts {
		levels[i], _ = strconv.Atoi(part)
	}
	return levels

}

func isSafe(levels []int) bool {
	increasing := levels[1] > levels[0]
	for i := 0; i < len(levels)-1; i++ {
		if increasing && levels[i] > levels[i+1] {
			return false
		}
		if !increasing && levels[i] < levels[i+1] {
			return false
		}
		d := abs(levels[i] - levels[i+1])
		if d < 1 {
			return false
		}
		if d > 3 {
			return false
		}
	}
	return true
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func allPossibleLevels(levels []int) iter.Seq[[]int] {
	return func(yield func([]int) bool) {
		if !yield(levels) {
			return
		}

		for i := 0; i < len(levels); i++ {
			l := []int{}
			l = append(l, levels[:i]...)
			l = append(l, levels[i+1:]...)
			if !yield(l) {
				return
			}
		}
	}
}

func anySafe(levels []int) bool {
	for l := range allPossibleLevels(levels) {
		if isSafe(l) {
			return true
		}
	}
	return false
}

func part2(input iter.Seq[string]) {
	s := 0
	for line := range input {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if anySafe(parseLine(line)) {
			s++
		}
	}
	fmt.Printf("Part 2: %d\n", s)
}
