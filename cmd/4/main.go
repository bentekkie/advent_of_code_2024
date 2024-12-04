package main

import (
	"flag"
	"fmt"
	"iter"
	"regexp"
	"strings"

	"github.com/bentekkie/advent_of_code_2024/pkg/inputs"
)

func main() {
	flag.Parse()
	input := inputs.String()
	part1(input)
	part2(input)
}

var xmas = regexp.MustCompile("XMAS")
var samx = regexp.MustCompile("SAMX")

func checkLine(line string) int {
	return len(xmas.FindAllString(line, -1)) + len(samx.FindAllString(line, -1))
}

func part1(input string) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	rlines := make([]strings.Builder, len(lines[0]))
	sum := 0
	for _, line := range lines {
		for i, c := range line {
			rlines[i].WriteRune(c)
		}
		sum += checkLine(line)
	}
	for _, line := range rlines {
		sum += checkLine(line.String())
	}
	var d, dr strings.Builder
	for k := range len(lines) + len(lines[0]) - 1 {
		d.Reset()
		dr.Reset()
		y := k
		x := 0
		rx := len(lines[0]) - 1
		if y >= len(lines) {
			x = y - len(lines) + 1
			rx -= x
			y = len(lines) - 1
		}
		for ; x < len(lines[0]) && y >= 0; x, rx, y = x+1, rx-1, y-1 {
			d.WriteByte(lines[y][x])
			dr.WriteByte(lines[y][rx])
		}
		sum += checkLine(d.String()) + checkLine(dr.String())
	}
	fmt.Printf("Part 1: %d\n", sum)
}

func check3x3MAS(lines []string) bool {
	return lines[1][1] == 'A' &&
		((lines[0][0] == 'M' && lines[2][2] == 'S') ||
			(lines[0][0] == 'S' && lines[2][2] == 'M')) &&
		((lines[0][2] == 'M' && lines[2][0] == 'S') ||
			(lines[0][2] == 'S' && lines[2][0] == 'M'))
}

func boxes(lines []string) iter.Seq[[]string] {
	return func(yield func([]string) bool) {
		for y := range len(lines) - 2 {
			for x := range len(lines[y]) - 2 {
				if !yield([]string{
					lines[y][x : x+3],
					lines[y+1][x : x+3],
					lines[y+2][x : x+3],
				}) {
					return
				}
			}
		}
	}
}

func part2(input string) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	c := 0
	for box := range boxes(lines) {
		if check3x3MAS(box) {
			c++
		}
	}
	fmt.Printf("Part 2: %d\n", c)
}
