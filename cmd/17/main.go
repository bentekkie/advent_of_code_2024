package main

import (
	"flag"
	"fmt"
	"iter"
	"maps"
	"math/big"
	"slices"
	"strconv"
	"strings"

	"github.com/bentekkie/advent_of_code_2024/pkg/inputs"
	"github.com/bentekkie/advent_of_code_2024/pkg/parse"
)

func main() {
	flag.Parse()
	input := inputs.String()
	part1(input)
	part2(input)
}

type Computer struct {
	regA, regB, regC *big.Int
	ip               int
	memory           []int
	output           []int
}

func (c *Computer) String() string {
	return fmt.Sprintf("A: %d\nB: %d\nC: %d\nIP: %d\nMemory: %v\nOutput: %v\n", c.regA, c.regB, c.regC, c.ip, c.memory, c.output)
}

func (c *Computer) combo(operand int) *big.Int {
	if operand < 4 {
		return big.NewInt(int64(operand))
	}
	if operand == 4 {
		return c.regA
	}
	if operand == 5 {
		return c.regB
	}
	if operand == 6 {
		return c.regC
	}
	return big.NewInt(0)
}

func (c *Computer) run() bool {
	steps := 0
	for c.ip < len(c.memory) {
		if !c.step() {
			return false
		}
		steps++
	}
	return true
}

var (
	two   = big.NewInt(2)
	eight = big.NewInt(8)
)

func (c *Computer) step() bool {
	code := c.memory[c.ip]
	operand := c.memory[c.ip+1]
	switch code {
	case 0: // adv
		c.regA.Div(c.regA, exp2big(c.combo(operand)))
	case 1: // bxl
		c.regB.Xor(c.regB, big.NewInt(int64(operand)))
	case 2: // bst
		c.regB.Mod(c.combo(operand), eight)
	case 3: // jnz
		if c.regA.Cmp(big.NewInt(0)) != 0 {
			c.ip = operand
			return true
		}
	case 4: // bxc
		c.regB.Xor(c.regB, c.regC)
	case 5: // out
		nextOut := new(big.Int).Mod(c.combo(operand), eight)
		c.output = append(c.output, int(nextOut.Int64()))
	case 6: // bdv
		c.regB.Div(c.regA, exp2big(c.combo(operand)))
	case 7: // cdv
		c.regC.Div(c.regA, exp2big(c.combo(operand)))
	}
	c.ip += 2
	return true
}

func exp2big(n *big.Int) *big.Int {
	return new(big.Int).Exp(two, n, nil)
}

func joinInts(nums []int) string {
	strs := make([]string, len(nums))
	for i, n := range nums {
		strs[i] = strconv.Itoa(n)
	}
	return strings.Join(strs, ",")
}

func part1(input string) {
	lines := strings.Split(input, "\n")
	a := parse.MustAtoi[int64](strings.TrimPrefix(lines[0], "Register A: "))
	b := parse.MustAtoi[int64](strings.TrimPrefix(lines[1], "Register B: "))
	c := parse.MustAtoi[int64](strings.TrimPrefix(lines[2], "Register C: "))
	mem := parse.NumList[int](strings.TrimPrefix(lines[4], "Program: "), ",")
	comp := &Computer{
		regA:   big.NewInt(a),
		regB:   big.NewInt(b),
		regC:   big.NewInt(c),
		memory: mem,
	}
	comp.run()
	fmt.Printf("Part 1: %s\n", joinInts(comp.output))
}

func getout(a *big.Int, mem []int) []int {
	comp := &Computer{
		regA:   a,
		regB:   big.NewInt(int64(0)),
		regC:   big.NewInt(int64(0)),
		memory: mem,
	}
	comp.run()
	return comp.output
}

func nextValid(a *big.Int, mem []int, want []int) iter.Seq[int64] {
	return func(yield func(int64) bool) {
		for n := int64(0); n < 8; n++ {
			testA := new(big.Int).Mul(a, eight)
			testA.Add(testA, big.NewInt(n))
			cpy := new(big.Int).Set(testA)
			out := getout(cpy, mem)
			if slices.Equal(out, want) {
				if yield(testA.Int64()) {
					return
				}
			}
		}
	}
}

func part2(input string) {
	lines := strings.Split(input, "\n")
	mem := parse.NumList[int](strings.TrimPrefix(lines[4], "Program: "), ",")
	nss := []map[int64]struct{}{{0: {}}}
	for k := range mem {
		ns := map[int64]struct{}{}
		for n := range nss[len(nss)-1] {
			for next := range nextValid(big.NewInt(n), mem, mem[len(mem)-1-k:]) {
				ns[next] = struct{}{}
			}
		}
		nss = append(nss, ns)
	}
	fmt.Printf("Part 2: %v\n", min(slices.Collect(maps.Keys(nss[len(nss)-1]))...))
}

func min[T ~int64](a ...T) T {
	min := a[0]
	for _, v := range a[1:] {
		if v < min {
			min = v
		}
	}
	return min
}
