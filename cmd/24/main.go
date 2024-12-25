package main

import (
	_ "embed"
	"flag"
	"fmt"
	"iter"
	"maps"
	"slices"
	"strings"

	"github.com/bentekkie/advent_of_code_2024/pkg/inputs"
	"github.com/bentekkie/advent_of_code_2024/pkg/parse"
)

var graphOut = flag.String("out", "", "")

func main() {
	flag.Parse()
	input := inputs.String()
	part1(input)
	part2(input)
}

type Gate interface {
	SetInputs(a, b Gate)
	GetInputs() (Gate, Gate)
	Set() bool
	Get() (val bool, ok bool)
	Clear()
	IsSet() bool
}

type Base struct {
	a, b     Gate
	set, val bool
}

type InitialWire struct {
	val bool
}

func (g *InitialWire) Clear() {}

func (g *InitialWire) Set() bool {
	return true
}

func (g *InitialWire) Get() (bool, bool) {
	return g.val, true
}

func (g *InitialWire) IsSet() bool {
	return true
}

func (g *InitialWire) SetInputs(a, b Gate) {
	panic("ahhhhh")
}

func (g *InitialWire) GetInputs() (Gate, Gate) {
	panic("ahhhhh")
}

func (g *Base) Clear() {
	g.set = false
}

func (g *Base) SetInputs(a, b Gate) {
	if a == nil || b == nil {
		panic("ahhhhh")
	}
	g.a, g.b = a, b
}

func (g *Base) GetInputs() (Gate, Gate) {
	return g.a, g.b
}

func (g *Base) IsSet() bool {
	return g.set
}

func (g *Base) Get() (val bool, ok bool) {
	return g.val, g.set
}

func (g *Base) GetInputVals() (a bool, b bool, ok bool) {
	a, ok = g.a.Get()
	if !ok {
		return false, false, false
	}
	b, ok = g.b.Get()
	if !ok {
		return false, false, false
	}
	return
}

type And struct {
	Base
}

func (g *And) Set() bool {
	a, b, ok := g.GetInputVals()
	if !ok {
		return false
	}
	g.val = a && b
	g.set = true
	return true
}

type Or struct {
	Base
}

func (g *Or) Set() bool {
	a, b, ok := g.GetInputVals()
	if !ok {
		return false
	}
	g.val = a || b
	g.set = true
	return true
}

type XOr struct {
	Base
}

func (g *XOr) Set() bool {
	a, b, ok := g.GetInputVals()
	if !ok {
		return false
	}
	g.val = a != b
	g.set = true
	return true
}

func iter2d[T any](it iter.Seq[T]) iter.Seq2[T, T] {
	return func(yield func(T, T) bool) {
		prev := []T{}
		for v1 := range it {
			for _, v2 := range prev {
				if !yield(v1, v2) {
					return
				}
			}
			prev = append(prev, v1)
		}
	}
}

func part1(input string) {
	parts := strings.Split(strings.TrimSpace(input), "\n\n")
	gates := map[string]Gate{}
	inputs := map[string][]string{}
	for _, initialStr := range strings.Split(strings.TrimSpace(parts[0]), "\n") {
		parts := strings.Split(strings.TrimSpace(initialStr), ": ")
		w := &InitialWire{val: parts[1] == "1"}
		gates[parts[0]] = w
	}

	for _, gateStr := range strings.Split(strings.TrimSpace(parts[1]), "\n") {
		gateStr = strings.TrimSpace(gateStr)
		if gateStr == "" {
			continue
		}
		parts := strings.Split(gateStr, " ")
		switch parts[1] {
		case "AND":
			gates[parts[4]] = &And{}
		case "OR":
			gates[parts[4]] = &Or{}
		case "XOR":
			gates[parts[4]] = &XOr{}
		}
		inputs[parts[4]] = []string{parts[0], parts[2]}
	}
	var toSet []string
	for gateName, inputNames := range inputs {
		gates[gateName].SetInputs(gates[inputNames[0]], gates[inputNames[1]])
		toSet = append(toSet, gateName)
	}
	for len(toSet) > 0 {
		newToSet := make([]string, 0, len(toSet))
		for _, g := range toSet {
			if !gates[g].Set() {
				newToSet = append(newToSet, g)
			}
		}
		toSet = newToSet
	}
	bits := map[int]bool{}
	for gateName, gate := range gates {
		if gateName[0] == 'z' {
			v, ok := gate.Get()
			if !ok {
				panic("ahhhhsjs")
			}
			bits[parse.MustAtoi[int](gateName[1:])] = v
		}
	}

	s := 0
	for bitNum, set := range bits {
		if set {
			s += 1 << bitNum
		}
	}
	fmt.Printf("Part 1: %v\n", s)
}

func part2(input string) {
	parts := strings.Split(strings.TrimSpace(input), "\n\n")
	infos := map[string]gateInfo{}
	for _, gateStr := range strings.Split(strings.TrimSpace(parts[1]), "\n") {
		gateStr = strings.TrimSpace(gateStr)
		if gateStr == "" {
			continue
		}
		parts := strings.Split(gateStr, " ")
		a, b := parts[0], parts[2]
		if b < a {
			a, b = b, a
		}
		infos[parts[4]] = gateInfo{t: parts[1], a: a, b: b}
	}
	correct := 1
	swaps := [][]string{}
	k := 0
	for correct < 45 {
		checkCorrect, usedTotal := correctUpTo(infos)
		correct = checkCorrect
		if correct < 45 {
			for s1, s2 := range iter2d(usedTotal.InOrder()) {
				swap(infos, s1, s2)
				k++
				testCorrect, _ := correctUpTo(infos)
				if testCorrect > correct {
					correct = testCorrect
					swaps = append(swaps, []string{s1, s2})
					break
				}
				swap(infos, s1, s2)
			}
		}
	}
	var names []string
	for _, swap := range swaps {
		names = append(names, swap...)
	}
	slices.Sort(names)
	fmt.Printf("Part 2: %d %s\n", k, strings.Join(names, ","))
}

func swap(infos map[string]gateInfo, a, b string) {
	infos[a], infos[b] = infos[b], infos[a]
}

func genLbls(prefix string, n int) []string {
	var lbls []string
	for i := 0; i < n; i++ {
		lbls = append(lbls, fmt.Sprintf("%s%02d", prefix, i))
	}
	return lbls
}

var (
	xs = genLbls("x", 50)
	ys = genLbls("y", 50)
)

func correctUpTo(infos map[string]gateInfo) (int, usedMap) {
	cin, _ := findWithInputs(xs[0], ys[0], "AND", infos)
	used := make(map[string]int, len(infos))
	for bit := 1; bit < 45; bit++ {
		cout, err := findCout(xs[bit], ys[bit], cin, infos, used)
		if err != nil {
			return bit, used
		}
		cin = cout
	}
	return 45, used
}

type gateInfo struct {
	t    string
	a, b string
}

func findWithInputs(a, b, t string, gates map[string]gateInfo) (string, bool) {
	for n, info := range gates {
		if info.t == t && info.a == a && info.b == b {
			return n, true
		}
		if info.t == t && info.a == b && info.b == a {
			return n, true
		}
	}
	return "", false
}

type usedMap map[string]int

func (u usedMap) Add(s string) bool {
	if _, ok := u[s]; ok {
		return false
	}
	u[s] = len(u)
	return true
}

func (u usedMap) InOrder() iter.Seq[string] {
	return func(yield func(string) bool) {
		byVal := map[int]string{}
		for k, v := range u {
			byVal[v] = k
		}
		keys := slices.Sorted(maps.Keys(byVal))
		slices.Reverse(keys)
		for _, k := range keys {
			if !yield(byVal[k]) {
				return
			}
		}
	}
}

func findCout(x, y, cin string, gates map[string]gateInfo, used usedMap) (string, error) {
	xor1, found := findWithInputs(x, y, "XOR", gates)
	if !found {
		return "", fmt.Errorf("unable to find xor1")
	}
	used.Add(xor1)
	and1, found := findWithInputs(x, y, "AND", gates)
	if !found {
		return "", fmt.Errorf("unable to find and1")
	}
	used.Add(and1)
	out, found := findWithInputs(cin, xor1, "XOR", gates)
	if !found {
		return "", fmt.Errorf("unable to find out")
	}
	used.Add(out)
	and2, found := findWithInputs(cin, xor1, "AND", gates)
	if !found {
		return "", fmt.Errorf("unable to find and2")
	}
	used.Add(and2)
	cout, found := findWithInputs(and1, and2, "OR", gates)
	if !found {
		return "", fmt.Errorf("unable to find cout")
	}
	used.Add(cout)
	return cout, nil
}
