package main

import (
	_ "embed"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

//go:embed input.txt
var inputtxt string

//go:embed example.txt
var exampletxt string

func main() {
	input := inputtxt
	part1(input)
	part2(input)
}

func part1(input string) {
	lines := strings.Split(input, "\n")
	var left []int
	var right []int
	for _, line := range lines {
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

func part2(input string) {
	lines := strings.Split(input, "\n")
	var left []int
	var right []int
	for _, line := range lines {
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
