package main

import (
	"flag"
	"fmt"
	"regexp"
	"strconv"
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
	r, _ := regexp.Compile(`mul\((\d+),(\d+)\)`)
	res := r.FindAllStringSubmatch(input, -1)
	sum := 0
	for _, match := range res {
		a, _ := strconv.Atoi(match[1])
		b, _ := strconv.Atoi(match[2])
		sum += a * b
	}
	fmt.Printf("Part 1: %v\n", sum)
}

func part2(input string) {
	r, _ := regexp.Compile(`mul\((\d+),(\d+)\)|do\(\)|don't\(\)`)
	res := r.FindAllStringSubmatch(input, -1)
	sum := 0
	enabled := true
	for _, match := range res {
		if match[0] == "do()" {
			enabled = true
		}
		if match[0] == "don't()" {
			enabled = false
		}
		if strings.HasPrefix(match[0], "mul") && enabled {
			a, _ := strconv.Atoi(match[1])
			b, _ := strconv.Atoi(match[2])
			sum += a * b
		}
	}
	fmt.Printf("Part 2: %v\n", sum)
}
