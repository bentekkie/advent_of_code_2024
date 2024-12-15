package main

import (
	"flag"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/bentekkie/advent_of_code_2024/pkg/inputs"
)

func main() {
	flag.Parse()
	input := inputs.String()
	part1(input)
	part2(input)
}

type Point struct {
	x, y int
}

func (p Point) add(o Point) Point {
	return Point{p.x + o.x, p.y + o.y}
}

type Box struct {
	ps []Point
}

type Grid struct {
	walls map[Point]struct{}
	boxes map[Point]struct{}
	robot Point
}

type WideGrid struct {
	walls map[Point]struct{}
	boxes map[Point]*Box
	robot Point
}

func (g *Grid) isBox(p Point) bool {
	_, ok := g.boxes[p]
	return ok
}

func (g *Grid) isWall(p Point) bool {
	_, ok := g.walls[p]
	return ok
}

func parse(input string) (*Grid, string) {
	parts := strings.Split(input, "\n\n")
	g := &Grid{
		walls: make(map[Point]struct{}),
		boxes: make(map[Point]struct{}),
	}
	for y, line := range strings.Split(strings.TrimSpace(parts[0]), "\n") {
		line = strings.TrimSpace(line)
		for x, c := range line {
			switch c {
			case '#':
				g.walls[Point{x, y}] = struct{}{}
			case '@':
				g.robot = Point{x, y}
			case 'O':
				g.boxes[Point{x, y}] = struct{}{}
			}
		}
	}

	return g, strings.ReplaceAll(parts[1], "\n", "")
}

func parseExpand(input string) (*WideGrid, string) {
	parts := strings.Split(input, "\n\n")
	g := &WideGrid{
		walls: make(map[Point]struct{}),
		boxes: make(map[Point]*Box),
	}
	for y, line := range strings.Split(strings.TrimSpace(parts[0]), "\n") {
		line = strings.TrimSpace(line)
		for x, c := range line {
			switch c {
			case '#':
				g.walls[Point{2 * x, y}] = struct{}{}
				g.walls[Point{2*x + 1, y}] = struct{}{}
			case '@':
				g.robot = Point{2 * x, y}
			case 'O':
				b := &Box{ps: []Point{{2 * x, y}, {2*x + 1, y}}}
				g.boxes[Point{2 * x, y}] = b
				g.boxes[Point{2*x + 1, y}] = b
			}
		}
	}

	return g, strings.ReplaceAll(parts[1], "\n", "")
}

func (g *WideGrid) boxAt(p Point) *Box {
	return g.boxes[p]
}

func (g *WideGrid) isWall(p Point) bool {
	_, ok := g.walls[p]
	return ok
}

func part1(input string) {
	g, inst := parse(input)
	for _, c := range inst {
		switch c {
		case '<':
			next := Point{g.robot.x - 1, g.robot.y}
			if g.isWall(next) {
				continue
			}
			if g.isBox(next) {
				empty := next
				for g.isBox(empty) {
					empty.x--
				}
				if g.isWall(empty) {
					continue
				}
				delete(g.boxes, next)
				g.boxes[empty] = struct{}{}
				g.robot.x--
			} else {
				g.robot.x--
			}
		case '>':
			next := Point{g.robot.x + 1, g.robot.y}
			if g.isWall(next) {
				continue
			}
			if g.isBox(next) {
				empty := next
				for g.isBox(empty) {
					empty.x++
				}
				if g.isWall(empty) {
					continue
				}
				delete(g.boxes, next)
				g.boxes[empty] = struct{}{}
				g.robot.x++
			} else {
				g.robot.x++
			}
		case '^':
			next := Point{g.robot.x, g.robot.y - 1}
			if g.isWall(next) {
				continue
			}
			if g.isBox(next) {
				empty := next
				for g.isBox(empty) {
					empty.y--
				}
				if g.isWall(empty) {
					continue
				}
				delete(g.boxes, next)
				g.boxes[empty] = struct{}{}
				g.robot.y--
			} else {
				g.robot.y--
			}
		case 'v':
			next := Point{g.robot.x, g.robot.y + 1}
			if g.isWall(next) {
				continue
			}
			if g.isBox(next) {
				empty := next
				for g.isBox(empty) {
					empty.y++
				}
				if g.isWall(empty) {
					continue
				}
				delete(g.boxes, next)
				g.boxes[empty] = struct{}{}
				g.robot.y++
			} else {
				g.robot.y++
			}
		}
	}
	s := 0
	for b := range g.boxes {
		s += b.x + 100*b.y
	}
	fmt.Printf("Part 1: %d\n", s)
}

func (g *WideGrid) moveRobot(delta Point) {
	next := g.robot.add(delta)
	if g.isWall(next) {
		return
	}
	nextBox := g.boxAt(next)
	if nextBox == nil {
		g.robot = next
		return
	}
	boxes := map[*Box]struct{}{g.boxAt(next): {}}
	toCheck := g.boxAt(next).ps
	checked := map[Point]struct{}{}
	for len(toCheck) > 0 {
		newToCheck := map[Point]struct{}{}
		for _, p := range toCheck {
			if _, ok := checked[p]; ok {
				continue
			}
			checked[p] = struct{}{}
			pnext := p.add(delta)
			if g.isWall(pnext) {
				return
			}
			if b := g.boxAt(pnext); b != nil {
				boxes[b] = struct{}{}
				for _, p := range b.ps {
					newToCheck[p] = struct{}{}
				}
			}
		}
		toCheck = slices.Collect(maps.Keys(newToCheck))
	}
	for b := range boxes {
		for i := range b.ps {
			delete(g.boxes, b.ps[i])
			b.ps[i] = b.ps[i].add(delta)
		}
	}
	for b := range boxes {
		for i := range b.ps {
			g.boxes[b.ps[i]] = b
		}
	}
	g.robot = next
}

var deltas = map[rune]Point{
	'<': {-1, 0},
	'>': {1, 0},
	'^': {0, -1},
	'v': {0, 1},
}

func part2(input string) {
	g, inst := parseExpand(input)
	for _, c := range inst {
		g.moveRobot(deltas[c])
	}
	s := 0
	boxes := map[Point]struct{}{}
	for _, b := range g.boxes {
		boxes[b.ps[0]] = struct{}{}
	}
	for b := range boxes {
		s += b.x + 100*b.y
	}
	fmt.Printf("Part 2: %d\n", s)
}
