package main

import (
	_ "embed"
	"flag"
	"fmt"
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
	canSwap := slices.Collect(maps.Keys(infos))
	usedCorrect := map[string]struct{}{}
	for correct < 45 {
		checkCorrect, usedTotal := correctUpTo(infos)
		correct = checkCorrect
		if correct < 45 {
			found := false
			for s1 := range usedTotal {
				if _, ok := usedCorrect[s1]; ok {
					continue
				}
				if found {
					break
				}
				for _, s2 := range canSwap {
					if _, ok := usedCorrect[s2]; ok {
						continue
					}
					swap(infos, s1, s2)
					testCorrect, _ := correctUpTo(infos)
					if testCorrect > correct {
						correct = testCorrect
						swaps = append(swaps, []string{s1, s2})
						found = true
						usedCorrect = usedTotal
						break
					}
					swap(infos, s1, s2)

				}
			}
			if !found {
				panic("aaaaa")
			}
		}
	}
	var names []string
	for _, swap := range swaps {
		names = append(names, swap...)
	}
	slices.Sort(names)
	fmt.Printf("Part 2: %s\n", strings.Join(names, ","))
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

func correctUpTo(infos map[string]gateInfo) (int, map[string]struct{}) {
	cin, _ := findWithInputs(xs[0], ys[0], "AND", infos)
	used := make(map[string]struct{}, len(infos))
	for bit := 1; bit < 45; bit++ {
		cout, _, usedStep, err := findCout(xs[bit], ys[bit], cin, infos)
		for _, name := range usedStep {
			used[name] = struct{}{}
		}
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

func findCout(x, y, cin string, gates map[string]gateInfo) (string, string, []string, error) {
	used := []string{x, y, cin}
	xor1, found := findWithInputs(x, y, "XOR", gates)
	if !found {
		return "", "", used, fmt.Errorf("unable to find xor1")
	}
	used = append(used, xor1)
	and1, found := findWithInputs(x, y, "AND", gates)
	if !found {
		return "", "", used, fmt.Errorf("unable to find and1")
	}
	used = append(used, and1)
	out, found := findWithInputs(cin, xor1, "XOR", gates)
	if !found {
		return "", "", used, fmt.Errorf("unable to find out")
	}
	used = append(used, and1)
	and2, found := findWithInputs(cin, xor1, "AND", gates)
	if !found {
		return "", "", used, fmt.Errorf("unable to find and2")
	}
	used = append(used, and2)
	cout, found := findWithInputs(and1, and2, "OR", gates)
	if !found {
		return "", "", used, fmt.Errorf("unable to find cout")
	}
	used = append(used, cout)
	return cout, out, used, nil
}
