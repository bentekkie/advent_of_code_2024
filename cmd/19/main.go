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

func part1(input string) {
	parts := strings.Split(strings.TrimSpace(input), "\n\n")
	towels := strings.Split(strings.TrimSpace(parts[0]), ", ")
	paterns := strings.Split(strings.TrimSpace(parts[1]), "\n")
	s := 0
	seen := make(map[string]int)
	for _, pattern := range paterns {
		if checkPattern(pattern, towels, seen) > 0 {
			s++
		}
	}
	fmt.Printf("Part 1: %v\n", s)
}

func checkPattern(pattern string, towels []string, seen map[string]int) int {
	if pattern == "" {
		return 1
	}
	if v, ok := seen[pattern]; ok {
		return v
	}
	s := 0
	for _, towel := range towels {
		if strings.HasPrefix(pattern, towel) {
			s += checkPattern(pattern[len(towel):], towels, seen)
		}
	}
	seen[pattern] = s
	return s
}

func part2(input string) {
	parts := strings.Split(strings.TrimSpace(input), "\n\n")
	towels := strings.Split(strings.TrimSpace(parts[0]), ", ")
	paterns := strings.Split(strings.TrimSpace(parts[1]), "\n")
	s := 0
	seen := make(map[string]int)
	for _, pattern := range paterns {
		s += checkPattern(pattern, towels, seen)
	}
	fmt.Printf("Part 2: %v\n", s)
}
