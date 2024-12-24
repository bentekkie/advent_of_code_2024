package main

import (
	"flag"
	"fmt"
	"iter"
	"maps"
	"math"
	"slices"
	"strings"

	"github.com/bentekkie/advent_of_code_2024/pkg/bengraph"
	"github.com/bentekkie/advent_of_code_2024/pkg/inputs"
	"github.com/bentekkie/advent_of_code_2024/pkg/parse"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/simple"
)

func main() {
	flag.Parse()
	part1(inputs.Lines())
	part2(inputs.Lines())
}

var pp = addMap(padPaths(dirKeys), padPaths(numpadKeys))

func addMap[K comparable, V any](ms ...map[K]V) map[K]V {
	totalLen := 0
	for _, m := range ms {
		totalLen += len(m)
	}
	m := make(map[K]V, totalLen)
	for _, mm := range ms {
		maps.Copy(m, mm)
	}
	return m
}

func part1(input iter.Seq[string]) {
	s := 0
	mem := map[expandState]int{}
	for line := range input {
		line := strings.TrimSpace(line)
		if line == "" {
			continue
		}
		numPart := parse.MustAtoi[int](line[:len(line)-1])
		short := expandCode(line, 3, mem)
		s += short * numPart
		fmt.Printf("DEBUG mem size: %d\n", len(mem))
	}
	fmt.Printf("Part 1: %v\n", s)
}

type expandState struct {
	code  string
	level int8
}

func expandMove(m Move, level int8, mem map[expandState]int) int {
	if level == 1 {
		return len(pp[m][0])
	}
	min := math.MaxInt
	for _, p := range pp[m] {
		if nextMin := expandCode(p, level-1, mem); nextMin < min {
			min = nextMin
		}
	}
	return min
}

func expandCode(p string, level int8, mem map[expandState]int) int {
	s := expandState{level: level, code: p}
	if ret, ok := mem[s]; ok {
		return ret
	}
	min := 0
	curr := 'A'
	for _, c := range p {
		min += expandMove(Move{from: runeidx[curr], to: runeidx[c]}, level, mem)
		curr = c
	}
	mem[s] = min
	return mem[s]
}

func part2(input iter.Seq[string]) {
	s := 0
	mem := map[expandState]int{}
	for line := range input {
		line := strings.TrimSpace(line)
		if line == "" {
			continue
		}
		numPart := parse.MustAtoi[int](line[:len(line)-1])
		short := expandCode(line, 26, mem)
		s += short * numPart
		fmt.Printf("DEBUG mem size: %d\n", len(mem))
	}
	fmt.Printf("Part 2: %v\n", s)
}

/*
+---+---+---+
| 7 | 8 | 9 |
+---+---+---+
| 4 | 5 | 6 |
+---+---+---+
| 1 | 2 | 3 |
+---+---+---+

	| 0 | A |
	+---+---+
*/

type runeid byte

func genRuneIDX(ms ...iter.Seq[rune]) map[rune]runeid {
	idx := map[rune]runeid{}
	for _, m := range ms {
		for r := range m {
			if _, ok := idx[r]; ok {
				continue
			}
			idx[r] = runeid(len(idx))
		}
	}
	return idx
}

var runeidx = genRuneIDX(maps.Keys(numpadKeys), maps.Keys(dirKeys))

var numpadKeys = map[rune]complex128{
	'7': 0 + 0i, '8': 1 + 0i, '9': 2 + 0i,
	'4': 0 + 1i, '5': 1 + 1i, '6': 2 + 1i,
	'1': 0 + 2i, '2': 1 + 2i, '3': 2 + 2i,
	/*	      */ '0': 1 + 3i, 'A': 2 + 3i,
}

/*
    +---+---+
    | ^ | A |
+---+---+---+
| < | v | > |
+---+---+---+
*/

var dirKeys = map[rune]complex128{
	/*	      */ '^': 1 + 0i, 'A': 2 + 0i,
	'<': 0 + 1i, 'v': 1 + 1i, '>': 2 + 1i,
}

var dirs = [4]complex128{1, -1, 1i, -1i}

type Move struct {
	from, to runeid
}

func (m Move) String() string {
	return fmt.Sprintf("'%c'->'%c'", m.from, m.to)
}

type bounds struct {
	max, min int8
}

func padPaths(keys map[rune]complex128) map[Move][]string {
	g := simple.NewUndirectedGraph()
	nodesByPos := map[complex128]graph.Node{}
	nodesByKey := map[rune]graph.Node{}
	for k, pos := range keys {
		n := bengraph.NewNode(g.NewNode(), pos)
		g.AddNode(n)
		nodesByPos[pos] = n
		nodesByKey[k] = n
	}
	for pos, n := range nodesByPos {
		for _, d := range dirs {
			next := pos + d
			if nextNode, ok := nodesByPos[next]; ok {
				g.SetEdge(g.NewEdge(n, nextNode))
			}
		}
	}
	paths := path.DijkstraAllPaths(g)
	pathStrs := map[Move][]string{}
	for from := range keys {
		for to := range keys {
			if from == to {
				pathStrs[Move{runeidx[from], runeidx[to]}] = []string{"A"}
				continue
			}
			ps, _ := paths.AllBetween(nodesByKey[from].ID(), nodesByKey[to].ID())
			strPs := make([]string, 0, len(ps))
			for _, p := range ps {
				var sb strings.Builder
				curr := p[0].(*bengraph.Node[complex128]).Data
				for _, n := range p[1:] {
					next := n.(*bengraph.Node[complex128]).Data
					switch next - curr {
					case 1:
						sb.WriteRune('>')
					case -1:
						sb.WriteRune('<')
					case 1i:
						sb.WriteRune('v')
					case -1i:
						sb.WriteRune('^')
					}
					curr = next
				}
				sb.WriteRune('A')
				strPs = append(strPs, sb.String())
			}
			repeats := map[string]bounds{}
			for _, ps := range strPs {
				cnt := int8(1)
				curr := rune(ps[0])
				var maxRepeat int8
				minRepeat := int8(math.MaxInt8)
				for _, c := range ps[1 : len(ps)-1] {
					if c == curr {
						cnt++
					} else {
						if cnt > maxRepeat {
							maxRepeat = cnt
						}
						if cnt < minRepeat {
							minRepeat = cnt
						}
						cnt = 1
					}
					curr = c
				}
				if cnt > maxRepeat {
					maxRepeat = cnt
				}
				if cnt < minRepeat {
					minRepeat = cnt
				}
				repeats[ps] = bounds{
					max: maxRepeat,
					min: minRepeat,
				}
			}
			var maxRepeat int8
			for _, b := range repeats {
				if b.max > maxRepeat {
					maxRepeat = b.max
				}
			}
			for key, b := range repeats {
				if b.max != maxRepeat {
					delete(repeats, key)
				}
			}
			var maxMin int8
			for _, b := range repeats {
				if b.min > maxMin {
					maxMin = b.min
				}
			}
			for key, b := range repeats {
				if b.min != maxMin {
					delete(repeats, key)
				}
			}
			pathStrs[Move{runeidx[from], runeidx[to]}] = slices.Collect(maps.Keys(repeats))
		}
	}
	return pathStrs
}
