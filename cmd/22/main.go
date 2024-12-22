package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/bentekkie/advent_of_code_2024/pkg/flags"
	"github.com/bentekkie/advent_of_code_2024/pkg/inputs"
	"github.com/bentekkie/advent_of_code_2024/pkg/parse"
)

func main() {
	flag.Parse()
	defer flags.CPUProfile()()
	input := inputs.String()
	part1(input)
	part2(input)
}

func part1(input string) {
	nums := parse.NumList[uint64](strings.TrimSpace(input), "\n")
	for range 2000 {
		for i := range nums {
			nums[i] = gensecret(nums[i])
		}
	}
	s := uint64(0)
	for _, n := range nums {
		s += n
	}
	fmt.Printf("Part 1: %v\n", s)
}

func gensecret(num uint64) uint64 {
	num ^= num * 64
	num %= 16777216
	num ^= num / 32
	num %= 16777216
	num ^= num * 2048
	num %= 16777216
	return num
}

func part2(input string) {
	nums := parse.NumList[uint64](strings.TrimSpace(input), "\n")
	prices := make([][]int, len(nums))
	for i, n := range nums {
		prices[i] = []int{int(n % 10)}
	}
	for range 2000 {
		for i := range nums {
			nums[i] = gensecret(nums[i])
			prices[i] = append(prices[i], int(nums[i]%10))
		}
	}
	changes := make([][]int, 0, len(nums))
	for _, price := range prices {
		change := make([]int, 0, len(price))
		for i, p := range price[1:] {
			change = append(change, p-price[i])
		}
		changes = append(changes, change)
	}
	scores := make([]map[ChangeSeq]int, len(nums))
	for i, change := range changes {
		scores[i] = getScores(change, prices[i])
	}
	allChangeSeqs := make(map[ChangeSeq]int, 10000)
	for _, s := range scores {
		for cs, v := range s {
			allChangeSeqs[cs] += v
		}
	}
	max := 0
	for _, v := range allChangeSeqs {
		if v > max {
			max = v
		}
	}
	fmt.Printf("Part 2: %v\n", max)
}

type ChangeSeq struct {
	a, b, c, d int
}

func getScores(changes []int, prices []int) map[ChangeSeq]int {
	scores := make(map[ChangeSeq]int, len(changes))
	for i := range changes[4:] {
		cs := ChangeSeq{
			a: changes[i],
			b: changes[i+1],
			c: changes[i+2],
			d: changes[i+3],
		}
		if _, ok := scores[cs]; !ok {
			scores[cs] = prices[i+4]
		}
	}
	return scores
}
