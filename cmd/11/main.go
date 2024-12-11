package main

import (
	"flag"
	"fmt"
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

func part1(input string) {
	stones := strings.Split(strings.TrimSpace(input), " ")
	//fmt.Printf("%v\n", stones)
	for range 25 {
		newStones := []string{}
		for _, stone := range stones {
			if stone == "0" {
				newStones = append(newStones, "1")
			} else if len(stone)%2 == 0 {
				left, right := stone[:len(stone)/2], stone[len(stone)/2:]
				for right[0] == '0' && len(right) > 1 {
					right = right[1:]
				}
				newStones = append(newStones, left, right)
			} else {
				num, _ := strconv.Atoi(stone)
				newStones = append(newStones, strconv.Itoa(num*2024))
			}
		}
		stones = newStones
		//fmt.Printf("%v\n", stones)
	}
	// TODO: implement me
	fmt.Printf("Part 1: %d\n", len(stones))
}

func part2(input string) {
	parts := strings.Split(strings.TrimSpace(input), " ")
	stones := map[string]int{}
	for _, p := range parts {
		stones[p]++
	}
	for range 75 {
		newStones := map[string]int{}
		for stone, count := range stones {
			if stone == "0" {
				newStones["1"] += count
			} else if len(stone)%2 == 0 {
				left, right := stone[:len(stone)/2], stone[len(stone)/2:]
				for right[0] == '0' && len(right) > 1 {
					right = right[1:]
				}
				newStones[left] += count
				newStones[right] += count
			} else {
				num, _ := strconv.Atoi(stone)
				newStones[strconv.Itoa(num*2024)] += count
			}
		}
		stones = newStones
		//fmt.Printf("%v\n", stones)
	}
	total := 0
	for _, c := range stones {
		total += c
	}
	// TODO: implement me
	fmt.Printf("Part 2: %d\n", total)
}
