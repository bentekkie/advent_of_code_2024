package main

import (
	"flag"
	"fmt"
	"iter"
	"math"
	"strconv"
	"strings"

	"github.com/bentekkie/advent_of_code_2024/pkg/benlog"
	"github.com/bentekkie/advent_of_code_2024/pkg/inputs"
)

func main() {
	flag.Parse()
	benlog.Timed(func() {
		part1(inputs.Lines())
	})
	benlog.Timed(func() {
		part1RTL(inputs.Lines())
	})
	benlog.Timed(func() {
		part2(inputs.Lines())
	})
	benlog.Timed(func() {
		part2RTL(inputs.Lines())
	})
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
	fmt.Printf("Part 1: %d\n", s)
}

func part1RTL(input iter.Seq[string]) {
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
		for i := range len(strNums) {
			n, _ := strconv.Atoi(strNums[len(strNums)-1-i])
			nums = append(nums, n)
		}
		if isPossibleRTL(nums, testNum) {
			s += testNum
		}
	}
	fmt.Printf("Part 1 RTL: %d\n", s)
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

func isPossibleRTL(nums []int, curr int) bool {
	if len(nums) == 0 {
		return curr == 0
	}
	if curr >= nums[0] && isPossibleRTL(nums[1:], curr-nums[0]) {
		return true
	}
	if curr%nums[0] == 0 && isPossibleRTL(nums[1:], curr/nums[0]) {
		return true
	}
	return false
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
	return isPossible2(testNum, next, curr*10*int(math.Pow10(int(math.Floor(math.Log10(float64(nums[0]))))))+nums[0])
}

func part2RTL(input iter.Seq[string]) {
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
		for i := range len(strNums) {
			n, _ := strconv.Atoi(strNums[len(strNums)-1-i])
			nums = append(nums, n)
		}
		if isPossible2RTL(nums, testNum) {
			s += testNum
		}
	}
	// TODO: implement me
	fmt.Printf("Part 2 RTL: %d\n", s)
}

func isPossible2RTL(nums []int, curr int) bool {
	if len(nums) == 0 {
		return curr == 0
	}
	if curr >= nums[0] && isPossible2RTL(nums[1:], curr-nums[0]) {
		return true
	}
	if curr%nums[0] == 0 && isPossible2RTL(nums[1:], curr/nums[0]) {
		return true
	}
	if curr == nums[0] {
		return isPossible2RTL(nums[1:], 0)
	}
	if curr >= nums[0] && curr >= 10 {
		ten := int(math.Pow(10, math.Floor(math.Log10(float64(nums[0]))))) * 10
		next := curr - nums[0]
		if next%ten == 0 {
			return isPossible2RTL(nums[1:], next/ten)
		}
	}
	return false
}
