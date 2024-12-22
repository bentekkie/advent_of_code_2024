package main

import (
	"flag"
	"fmt"
	"iter"
	"maps"
	"math"
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
	mem := map[expandState]expandStateRet{}
	for line := range input {
		line := strings.TrimSpace(line)
		if line == "" {
			continue
		}
		numPart := parse.MustAtoi[int](line[:len(line)-1])
		short := minForCode(line, 2, mem)
		s += short * numPart
		fmt.Printf("DEBUG mem size: %d\n", len(mem))
	}
	fmt.Printf("Part 1: %v\n", s)
}

func minForCode(code string, keypads int, mem map[expandState]expandStateRet) int {
	return expandCode(code, strings.Repeat("A", keypads+1), mem).minLen
}

type expandState struct {
	currs string
	press rune
	code  string
}

type expandStateRet struct {
	minLen int
	currs  string
}

func expandMove(s expandState, mem map[expandState]expandStateRet) expandStateRet {
	paths := pp[Move{from: rune(s.currs[0]), to: s.press}]
	if _, ok := mem[s]; !ok && len(s.currs) == 1 {
		mem[s] = expandStateRet{minLen: len(paths[0]), currs: string(s.press)}
	} else if !ok {
		ret := expandStateRet{minLen: math.MaxInt64}
		for _, p := range paths {
			if next := expandCode(p, s.currs[1:], mem); next.minLen < ret.minLen {
				ret.minLen = next.minLen
				ret.currs = next.currs
			}
		}
		ret.currs = string(s.press) + ret.currs
		mem[s] = ret
	}
	return mem[s]
}

func expandCode(p string, currs string, mem map[expandState]expandStateRet) expandStateRet {
	s := expandState{currs: currs, code: p}
	if ret, ok := mem[s]; ok {
		return ret
	}
	ret := expandStateRet{currs: currs}
	for _, c := range p {
		next := expandMove(expandState{currs: ret.currs, press: c}, mem)
		ret.minLen += next.minLen
		ret.currs = next.currs
	}
	mem[s] = ret
	return mem[s]
}

func part2(input iter.Seq[string]) {
	s := 0
	mem := map[expandState]expandStateRet{}
	for line := range input {
		line := strings.TrimSpace(line)
		if line == "" {
			continue
		}
		numPart := parse.MustAtoi[int](line[:len(line)-1])
		short := minForCode(line, 25, mem)
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
	from, to rune
}

func (m Move) String() string {
	return fmt.Sprintf("'%c'->'%c'", m.from, m.to)
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
				pathStrs[Move{from, to}] = []string{"A"}
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
			pathStrs[Move{from, to}] = strPs
		}
	}
	return pathStrs
}
