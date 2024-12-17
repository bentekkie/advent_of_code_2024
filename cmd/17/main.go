package main

import (
	"flag"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/bentekkie/advent_of_code_2024/pkg/inputs"
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

func mustParseInt(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return n
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
	a := mustParseInt(strings.TrimPrefix(lines[0], "Register A: "))
	b := mustParseInt(strings.TrimPrefix(lines[1], "Register B: "))
	c := mustParseInt(strings.TrimPrefix(lines[2], "Register C: "))
	memStr := strings.Split(strings.TrimPrefix(lines[4], "Program: "), ",")
	mem := make([]int, len(memStr))
	for i, s := range memStr {
		mem[i] = mustParseInt(s)
	}
	comp := &Computer{
		regA:   big.NewInt(int64(a)),
		regB:   big.NewInt(int64(b)),
		regC:   big.NewInt(int64(c)),
		memory: mem,
	}
	comp.run()
	fmt.Printf("Part 1: %s\n", joinInts(comp.output))
}

func check(a *big.Int, mem []int) (bool, int, int, *Computer) {
	comp := &Computer{
		regA:   a,
		regB:   big.NewInt(int64(0)),
		regC:   big.NewInt(int64(0)),
		memory: mem,
	}
	comp.run()
	for i, n := range mem {
		if i >= len(comp.output) || comp.output[i] != n {
			return false, i, len(comp.output), comp
		}
	}
	return true, len(comp.output), len(comp.output), comp
}

func findCorrectForLen(a *big.Int, iseen int, mem []int) (*big.Int, int) {
	max := big.NewInt(1 << 19)
	fmt.Printf("Starting at %v\n", a)
	for n := big.NewInt(0); n.Cmp(max) < 0; n.Add(n, big.NewInt(1)) {
		testA := new(big.Int).Set(n)
		testA.Lsh(testA, uint(a.BitLen()))
		testA.Add(testA, a)
		cpy := new(big.Int).Set(testA)
		ok, i, l, _ := check(cpy, mem)
		if ok {
			return testA, i
		}
		if i > iseen && l == i && i != 15 {
			return testA, i
		}
	}
	panic("err")
}

func part2(input string) {
	lines := strings.Split(input, "\n")
	memStr := strings.Split(strings.TrimPrefix(lines[4], "Program: "), ",")
	mem := make([]int, len(memStr))
	for i, s := range memStr {
		mem[i] = mustParseInt(s)
	}
	maxi := 0
	a := big.NewInt(7)
	for maxi < len(mem) {
		a, maxi = findCorrectForLen(a, maxi, mem)
		fmt.Printf("a=%v maxi=%v %s\n", a, maxi, bitstring(a))
	}
	fmt.Printf("Part 2: %v", a)

}

func bitstring(n *big.Int) string {
	var sb strings.Builder
	for i := range n.BitLen() {
		if n.Bit(n.BitLen()-1-i) == 1 {
			sb.WriteRune('1')
		} else {
			sb.WriteRune('0')
		}
	}
	return sb.String()
}
