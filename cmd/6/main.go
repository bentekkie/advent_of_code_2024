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

type P struct {
	x, y int
}

func (p P) up() P {
	return P{p.x, p.y - 1}
}

func (p P) down() P {
	return P{p.x, p.y + 1}
}

func (p P) left() P {
	return P{p.x - 1, p.y}
}

func (p P) right() P {
	return P{p.x + 1, p.y}
}
func (p P) String() string {
	return fmt.Sprintf("(%d,%d)", p.x, p.y)
}
func (p P) inbounds(minx, miny, maxx, maxy int) bool {
	return p.x >= minx && p.x <= maxx && p.y >= miny && p.y <= maxy
}

type D int

const (
	up    D = 0
	right D = 1
	down  D = 2
	left  D = 3
)

func part1(input string) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	var guard P
	guarDir := up
	obsticals := map[P]struct{}{}
	for y, line := range lines {
		line := strings.TrimSpace(line)
		for x, c := range line {
			if c == '#' {
				obsticals[P{x, y}] = struct{}{}
			}
			if c == '^' {
				guard = P{x, y}
			}
		}
	}
	//fmt.Printf("Guard: %s\n", guard)
	//fmt.Printf("Obsticals: %v\n", obsticals)
	visited := map[P]struct{}{}
	for guard.inbounds(0, 0, len(lines[0])-1, len(lines)-1) {
		visited[guard] = struct{}{}
		var nextP P
		switch guarDir {
		case up:
			nextP = guard.up()
		case down:
			nextP = guard.down()
		case left:
			nextP = guard.left()
		case right:
			nextP = guard.right()
		}
		if _, ok := obsticals[nextP]; ok {
			guarDir = (guarDir + 1) % 4
		} else {
			guard = nextP
		}
	}
	fmt.Printf("Part 1: %d\n", len(visited))
}

type V struct {
	p P
	d D
}

func isLoop(start V, obsticals map[P]struct{}, maxx, maxy int) bool {
	guard := start.p
	guarDir := start.d
	visited := map[V]struct{}{}
	for guard.inbounds(0, 0, maxx, maxy) {
		if _, ok := visited[V{guard, guarDir}]; ok {
			return true
		}
		visited[V{guard, guarDir}] = struct{}{}
		var nextP P
		switch guarDir {
		case up:
			nextP = guard.up()
		case down:
			nextP = guard.down()
		case left:
			nextP = guard.left()
		case right:
			nextP = guard.right()
		}
		if _, ok := obsticals[nextP]; ok {
			guarDir = (guarDir + 1) % 4
		} else {
			guard = nextP
		}
	}
	return false
}

func path(start V, obsticals map[P]struct{}, maxx, maxy int) map[P]struct{} {
	guard := start.p
	guarDir := start.d
	visited := map[P]struct{}{}
	for guard.inbounds(0, 0, maxx, maxy) {
		visited[guard] = struct{}{}
		var nextP P
		switch guarDir {
		case up:
			nextP = guard.up()
		case down:
			nextP = guard.down()
		case left:
			nextP = guard.left()
		case right:
			nextP = guard.right()
		}
		if _, ok := obsticals[nextP]; ok {
			guarDir = (guarDir + 1) % 4
		} else {
			guard = nextP
		}
	}
	return visited
}

func part2(input string) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	var guard P
	guarDir := up
	obsticals := map[P]struct{}{}
	possibleObsticals := []P{}
	for y, line := range lines {
		line := strings.TrimSpace(line)
		for x, c := range line {
			if c == '#' {
				obsticals[P{x, y}] = struct{}{}
			}
			if c == '^' {
				guard = P{x, y}
			}
			if c == '.' {
				possibleObsticals = append(possibleObsticals, P{x, y})
			}
		}
	}
	original := path(V{guard, guarDir}, obsticals, len(lines[0])-1, len(lines)-1)
	var s int
	for _, o := range possibleObsticals {
		if _, ok := original[o]; !ok {
			continue
		}
		obsticals[o] = struct{}{}
		if isLoop(V{guard, guarDir}, obsticals, len(lines[0])-1, len(lines)-1) {
			s++
		}
		delete(obsticals, o)
	}

	fmt.Printf("Part s: %d\n", s)
}
