package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/bentekkie/advent_of_code_2024/pkg/benlog"
	"github.com/bentekkie/advent_of_code_2024/pkg/inputs"
)

func main() {
	flag.Parse()
	input := inputs.String()
	part1(input)
	part2(input)
}

func checkLine(line string) int {
	c := 0
	for i := range len(line) - 3 {
		if line[i:i+4] == "XMAS" {
			c++
		} else if line[i:i+4] == "SAMX" {
			c++
		}
	}
	return c
}

func part1(input string) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	sum := 0
	for _, line := range lines {
		sum += checkLine(line)
	}
	rlines := make([]string, 0, len(lines[0]))
	for range len(lines[0]) {
		rlines = append(rlines, "")
	}
	for _, line := range lines {
		for i, c := range line {
			if strings.TrimSpace(string(c)) == "" {
				continue
			}
			rlines[i] += string(c)
		}
	}
	for _, line := range rlines {
		sum += checkLine(line)
	}
	diags := []string{}
	for k := range len(lines) + len(lines[0]) {
		y := k
		x := 0
		diag := ""
		if y >= len(lines) {
			x = y - len(lines)
			y = len(lines) - 1
		}
		for ; x < len(lines[0]) && y >= 0; x, y = x+1, y-1 {
			diag += string(lines[y][x])
		}
		if k == 0 || diag != diags[len(diags)-1] {
			diags = append(diags, diag)
		}
	}
	for _, line := range diags {
		sum += checkLine(line)
	}
	rdiags := []string{}
	for k := range len(lines) + len(lines[0]) {
		y := k
		x := len(lines[0]) - 1
		diag := ""
		if y >= len(lines) {
			x = len(lines[0]) - 1 - (y - len(lines))
			y = len(lines) - 1
		}
		for ; x >= 0 && y >= 0; x, y = x-1, y-1 {
			diag += string(lines[y][x])
		}
		if k == 0 || diag != rdiags[len(rdiags)-1] {
			rdiags = append(rdiags, diag)
		}
	}
	for _, line := range rdiags {
		sum += checkLine(line)
	}
	fmt.Printf("Part 1: %d\n", sum)
}

func check3x3MAS(lines []string) bool {
	return lines[1][1] == 'A' && ((lines[0][0] == 'M' && lines[2][2] == 'S') || (lines[0][0] == 'S' && lines[2][2] == 'M')) && ((lines[0][2] == 'M' && lines[2][0] == 'S') || (lines[0][2] == 'S' && lines[2][0] == 'M'))
}

func part2(input string) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	c := 0
	for y := range len(lines) - 2 {
		for x := range len(lines[y]) - 2 {
			box := []string{
				lines[y][x : x+3],
				lines[y+1][x : x+3],
				lines[y+2][x : x+3],
			}
			benlog.ExamplePrintf("%s\n\n", strings.Join(box, "\n"))
			if check3x3MAS(box) {
				c++
			}
		}
	}
	fmt.Printf("Part 2: %d\n", c)
}
