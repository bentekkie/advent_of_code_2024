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

func part1(input iter.Seq[string]) {
	s := 0
	for line := range input {
		line := strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Split(line, ": ")
		testNum, _ := strconv.Atoi(parts[0])
		strNums := strings.Split(parts[1], " ")
		nums := make([]int, 0, len(strNums))
		for _, snum := range strNums {
			n, _ := strconv.Atoi(snum)
			nums = append(nums, n)
		}
		if isPossible(testNum, nums[1:], nums[0]) {
			s += testNum
		}

	}
	// TODO: implement me
	fmt.Printf("Part 1: %d\n", s)
}

func isPossible(testNum int, nums []int, curr int) bool {
	if curr == testNum && len(nums) == 0 {
		return true
	}
	if curr > testNum {
		return false
	}
	if len(nums) == 0 {
		return false
	}
	next := nums[1:]
	if isPossible(testNum, next, curr*nums[0]) {
		return true
	}
	return isPossible(testNum, next, curr+nums[0])
}

func part2(input iter.Seq[string]) {
	s := 0
	for line := range input {
		line := strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Split(line, ": ")
		testNum, _ := strconv.Atoi(parts[0])
		strNums := strings.Split(parts[1], " ")
		nums := make([]int, 0, len(strNums))
		for _, snum := range strNums {
			n, _ := strconv.Atoi(snum)
			nums = append(nums, n)
		}
		if isPossible2(testNum, nums[1:], nums[0]) {
			s += testNum
		}
	}
	// TODO: implement me
	fmt.Printf("Part 2: %d\n", s)
}

func isPossible2(testNum int, nums []int, curr int) bool {
	if curr == testNum && len(nums) == 0 {
		return true
	}
	if curr > testNum {
		return false
	}
	if len(nums) == 0 {
		return false
	}
	next := nums[1:]
	if isPossible2(testNum, next, curr*nums[0]) {
		return true
	}
	if isPossible2(testNum, next, curr+nums[0]) {
		return true
	}
	n, _ := strconv.Atoi(strconv.Itoa(curr) + strconv.Itoa(nums[0]))
	return isPossible2(testNum, next, n)
}
