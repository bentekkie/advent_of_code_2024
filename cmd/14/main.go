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

type Point struct {
	x int64
	y int64
}

type Velocity struct {
	x int64
	y int64
}

func posMod(n, m int64) int64 {
	return (n%m + m) % m
}

func (p Point) quadrant(max Point) int {
	if p.x < (max.x-1)/2 {
		if p.y < (max.y-1)/2 {
			return 1
		} else if p.y > (max.y-1)/2 {
			return 2
		}
	} else if p.x > (max.x-1)/2 {
		if p.y < (max.y-1)/2 {
			return 3
		} else if p.y > (max.y-1)/2 {
			return 4
		}
	}
	return 0
}

func mustParseInt(s string) int64 {
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return n
}

func part1(input iter.Seq[string]) {
	robotMap := map[Point][]Velocity{}
	for line := range input {
		parts := strings.Split(strings.TrimSpace(line), " ")
		if len(parts) != 2 {
			continue
		}
		pparts := strings.Split(strings.TrimPrefix(parts[0], "p="), ",")
		vparts := strings.Split(strings.TrimPrefix(parts[1], "v="), ",")
		p := Point{mustParseInt(pparts[0]), mustParseInt(pparts[1])}
		v := Velocity{mustParseInt(vparts[0]), mustParseInt(vparts[1])}
		robotMap[p] = append(robotMap[p], v)
	}
	max := Point{101, 103}
	for range 100 {
		robotMap = moveRobots(robotMap, max)
	}
	quadCounts := map[int]int{}
	for p, vs := range robotMap {
		quadCounts[p.quadrant(max)] += len(vs)
	}
	fmt.Printf("Part 1: %v\n", quadCounts[1]*quadCounts[2]*quadCounts[3]*quadCounts[4])
}

func triangle(top Point, height int) iter.Seq[Point] {
	return func(yield func(Point) bool) {
		yield(top)
		for row := range height - 1 {
			for x := top.x - int64(row); x <= top.x+int64(row); x++ {
				if !yield(Point{x, top.y + int64(row)}) {
					return
				}
			}
		}
	}
}

func containsTriangle(locs map[Point][]Velocity, top Point, height int) bool {
	for p := range triangle(top, height) {
		if vs, ok := locs[p]; !ok || len(vs) == 0 {
			return false
		}
	}
	return true
}

func largestTriangle(robots map[Point][]Velocity) int {
	l := 0
	for top := range robots {
		height := l
		for containsTriangle(robots, top, height) {
			height++
		}
		l = height
	}
	return l
}

func grid(robots map[Point][]Velocity, max Point) string {
	var sb strings.Builder
	for y := range max.y {
		for x := range max.x {
			if len(robots[Point{x, y}]) > 0 {
				sb.WriteRune('x')
			} else {
				sb.WriteRune('.')
			}
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}

func moveRobots(robotMap map[Point][]Velocity, max Point) map[Point][]Velocity {
	nextRobotMap := make(map[Point][]Velocity, int(max.x*max.y))
	for p, vs := range robotMap {
		for _, v := range vs {
			nextP := Point{
				x: posMod(p.x+v.x, max.x),
				y: posMod(p.y+v.y, max.y),
			}
			nextRobotMap[nextP] = append(nextRobotMap[nextP], v)
		}
	}
	return nextRobotMap
}

func part2(input iter.Seq[string]) {
	robotMap := map[Point][]Velocity{}
	for line := range input {
		parts := strings.Split(strings.TrimSpace(line), " ")
		if len(parts) != 2 {
			continue
		}
		pparts := strings.Split(strings.TrimPrefix(parts[0], "p="), ",")
		vparts := strings.Split(strings.TrimPrefix(parts[1], "v="), ",")
		p := Point{mustParseInt(pparts[0]), mustParseInt(pparts[1])}
		v := Velocity{mustParseInt(vparts[0]), mustParseInt(vparts[1])}
		robotMap[p] = append(robotMap[p], v)
	}
	max := Point{101, 103}
	for k := range max.x * max.y {
		robotMap = moveRobots(robotMap, max)
		m := largestTriangle(robotMap)
		if m > 5 {
			fmt.Printf("%s\nPart 2: %d\n", grid(robotMap, max), k+1)
			return
		}
	}
}
