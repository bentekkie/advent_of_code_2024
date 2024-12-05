package main

import (
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

func parseInput(lines iter.Seq[string]) (map[int]map[int]struct{}, [][]int) {
	updates := [][]int{}
	addingRules := true
	after := map[int]map[int]struct{}{}
	for line := range lines {
		line := strings.TrimSpace(line)
		if line == "" {
			addingRules = false
			continue
		}
		if addingRules {
			parts := strings.Split(line, "|")
			left, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
			right, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
			if _, ok := after[left]; !ok {
				after[left] = map[int]struct{}{}
			}
			after[left][right] = struct{}{}
		} else {
			parts := strings.Split(line, ",")
			update := []int{}
			for _, part := range parts {
				p, _ := strconv.Atoi(strings.TrimSpace(part))
				update = append(update, p)
			}
			updates = append(updates, update)
		}
	}
	return after, updates
}

func part1(lines iter.Seq[string]) {
	after, updates := parseInput(lines)
	sum := 0
	sortFunc := func(a, b int) int {
		if _, ok := after[a][b]; ok {
			return -1
		} else if _, ok := after[b][a]; ok {
			return 1
		} else {
			return 0
		}
	}
	for _, update := range updates {
		if slices.IsSortedFunc(update, sortFunc) {
			sum += update[len(update)/2]
		}
	}
	fmt.Printf("Part 1: %d\n", sum)
}

func part2(lines iter.Seq[string]) {
	after, updates := parseInput(lines)
	sum := 0
	sortFunc := func(a, b int) int {
		if _, ok := after[a][b]; ok {
			return -1
		} else if _, ok := after[b][a]; ok {
			return 1
		} else {
			return 0
		}
	}
	for _, update := range updates {
		if !slices.IsSortedFunc(update, sortFunc) {
			slices.SortFunc(update, sortFunc)
			sum += update[len(update)/2]
		}
	}
	fmt.Printf("Part 2: %d\n", sum)
}
