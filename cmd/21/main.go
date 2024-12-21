package main

import (
	"flag"
	"fmt"
	"iter"
	"maps"
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

var keyPadPaths = padPaths(dirKeys)
var numpadPaths = padPaths(numpadKeys)

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
	}
	fmt.Printf("Part 1: %v\n", s)
}

func minForCode(code string, keypads int, mem map[expandState]expandStateRet) int {
	numpaths := slices.Collect(maps.Keys(typeOnPad(code, numpadPaths)))
	minLen := -1
	for _, numpath := range numpaths {
		newMin := typeOnKeyPadLevels(numpath, keypads, mem)
		if newMin < minLen || minLen == -1 {
			minLen = newMin
		}
	}
	return minLen
}

type expandState struct {
	level int
	currs string
	press rune
}

type expandStateRet struct {
	minLen int
	currs  string
}

func expandMove(currs string, press rune, levels int, mem map[expandState]expandStateRet) (int, string) {
	s := expandState{level: levels, currs: currs, press: press}
	paths := keyPadPaths[Move{from: rune(currs[0]), to: press}]
	if _, ok := mem[s]; !ok && levels == 1 {
		mem[s] = expandStateRet{minLen: len(paths[0]), currs: string(press)}
	} else if !ok {
		ret := expandStateRet{
			minLen: -1,
			currs:  "",
		}
		for _, p := range paths {
			nextLevelCurr := currs[1:]
			pMin := 0
			for _, c := range p {
				nextMin, nextCurr := expandMove(nextLevelCurr, c, levels-1, mem)
				pMin += nextMin
				nextLevelCurr = nextCurr
				if pMin > ret.minLen && ret.minLen >= 0 {
					break
				}
			}
			if pMin < ret.minLen || ret.minLen < 0 {
				ret.minLen = pMin
				ret.currs = string(press) + nextLevelCurr
			}
		}
		mem[s] = ret
	}
	return mem[s].minLen, mem[s].currs
}

func typeOnKeyPadLevels(code string, levels int, mem map[expandState]expandStateRet) int {
	var sb strings.Builder
	for range levels {
		sb.WriteRune('A')
	}
	currs := sb.String()
	min := 0
	for _, c := range code {
		newMin, newCurrs := expandMove(currs, c, levels, mem)
		currs = newCurrs
		min += newMin
	}
	return min
}

func typeOnPad(code string, padPaths map[Move][]string) map[string]struct{} {
	currButton := 'A'
	paths := map[string]struct{}{"": {}}
	for _, c := range code {
		nextPaths := map[string]struct{}{}
		newPaths := padPaths[Move{from: currButton, to: c}]
		for p := range paths {
			for _, np := range newPaths {
				nextPaths[p+np] = struct{}{}
			}
		}
		currButton = c
		paths = nextPaths
	}
	return paths
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
