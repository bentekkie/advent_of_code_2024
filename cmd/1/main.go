package main

import (
	_ "embed"
	"flag"
	"fmt"
	"iter"
	"slices"
	"strconv"
	"strings"

	"github.com/bentekkie/advent_of_code_2024/pkg/inputs"
)

func main() {
	flag.Parse()
	part1(inputs.Lines())
	part2(inputs.Lines())
}

func part1(lines iter.Seq[string]) {
	var left []int
	var right []int
	for line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, "   ")
		l, _ := strconv.Atoi(parts[0])
		r, _ := strconv.Atoi(parts[1])
		left = append(left, l)
		right = append(right, r)
	}
	slices.Sort(left)
	slices.Sort(right)
	dists := 0
	for i := 0; i < len(left); i++ {
		d := right[i] - left[i]
		if d < 0 {
			d = -d
		}
		dists += d
	}
	fmt.Printf("Part 1: %d\n", dists)

}

func part2(lines iter.Seq[string]) {
	var left []int
	var right []int
	for line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, "   ")
		l, _ := strconv.Atoi(parts[0])
		r, _ := strconv.Atoi(parts[1])
		left = append(left, l)
		right = append(right, r)
	}
	counts := map[int]int{}
	for _, num := range right {
		counts[num]++
	}
	sum := 0
	for _, num := range left {
		sum += num * counts[num]
	}
	fmt.Printf("Part 2: %d\n", sum)

}
