package main

import (
	"flag"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/bentekkie/advent_of_code_2024/pkg/benlog"
	"github.com/bentekkie/advent_of_code_2024/pkg/inputs"
)

func main() {
	flag.Parse()
	input := inputs.String()
	part1(input)
	benlog.Timed(func() {
		part2(input)
	})
	benlog.Timed(func() {
		part2strmul(input)
	})
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

func part2strmul(input string) {
	parts := strings.Split(strings.TrimSpace(input), " ")
	stones := map[int]int{}
	for _, p := range parts {
		n, _ := strconv.Atoi(p)
		stones[n]++
	}
	for range 75 {
		newStones := map[int]int{}
		for stone, count := range stones {
			if stone == 0 {
				newStones[1] += count
			} else if l, r, ok := split(stone); ok {
				newStones[l] += count
				newStones[r] += count
			} else {
				newStones[2024*stone] += count
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

func split(stone int) (int, int, bool) {
	dl := diglen(stone)
	if dl%2 == 1 {
		return 0, 0, false
	}
	ten := int(math.Pow10(dl / 2))
	return stone / ten, stone % ten, true
}

func diglen(a int) int {
	return int(math.Floor(math.Log10(float64(a)))) + 1
}

func strMul(a int, b string) int {
	total := 0
	place := 1
	for i := len(b) - 1; i >= 0; i-- {
		switch b[i] {
		case '0':
		case '1':
			total += place * a
		case '2':
			total += place * 2 * a
		case '3':
			total += place * 3 * a
		case '4':
			total += place * 4 * a
		case '5':
			total += place * 5 * a
		case '6':
			total += place * 6 * a
		case '7':
			total += place * 7 * a
		case '8':
			total += place * 8 * a
		case '9':
			total += place * 9 * a
		}
		place *= 10
	}
	return total
}
