package main

import (
	"flag"
	"fmt"
	"slices"
	"strconv"

	"github.com/bentekkie/advent_of_code_2024/pkg/benlog"
	"github.com/bentekkie/advent_of_code_2024/pkg/inputs"
)

func main() {
	flag.Parse()
	input := inputs.String()
	benlog.Timed(func() {
		part1(input)
	})
	benlog.Timed(func() {
		part2(input)
	})
}

type section struct {
	prev, next *section
	id         int
	count      int
}

func part1(input string) {
	sects := make(map[int]*section, len(input))
	var prev *section
	for i, c := range input {
		if c == '\n' {
			continue
		}
		num, _ := strconv.Atoi(string(c))
		if i%2 == 0 {
			s := &section{
				id:    (i / 2),
				count: num,
				prev:  prev,
			}
			if s.prev != nil {
				s.prev.next = s
			}
			prev = s
			sects[s.id] = s
		} else {
			s := &section{
				id:    -1,
				count: num,
				prev:  prev,
			}
			if s.prev != nil {
				s.prev.next = s
			}
			prev = s
			sects[s.id] = s
		}
	}
	left := sects[0]
	for left.id != -1 {
		left = left.next
	}
	right := prev
	for right.id != 0 {
		if right.id == -1 {
			right = right.prev
			continue
		}
		if right.count < left.count {
			oldRight := right
			oldRightNext := right.next
			oldRightPrev := right.prev
			newRight := right.prev
			for newRight.id == -1 {
				newRight = newRight.prev
			}
			oldLeftPrev := left.prev
			oldLeftPrev.next = oldRight
			oldRight.prev = oldLeftPrev
			left.prev = oldRight
			left.count -= oldRight.count
			oldRight.next = left
			oldRightPrev.next = oldRightNext
			if oldRightNext != nil {
				oldRightNext.prev = oldRightPrev
			}
			right = newRight
		} else if right.count > left.count {
			nextLeft := left.next
			for nextLeft.id != -1 {
				nextLeft = nextLeft.next
			}
			left.id = right.id
			right.count -= left.count
			left = nextLeft
		} else {
			newRight := right.prev
			for newRight.id == -1 {
				newRight = newRight.prev
			}
			nextLeft := left.next
			for nextLeft.id != -1 {
				nextLeft = nextLeft.next
			}
			left.id = right.id
			right.id = -1
			right = newRight
			left = nextLeft
		}
		for left.next != nil && left.next.id == -1 {
			leftNext := left.next
			left.count += leftNext.count
			left.next = left.next.next
			if left.next != nil {
				left.next.prev = left
			}
			left = leftNext
		}
		if left.next == nil {
			break
		}
	}
	pos := 0
	s := 0
	curr := sects[0]
	for curr.id != -1 {
		for range curr.count {
			s += int(curr.id) * pos
			pos++
		}
		curr = curr.next
	}
	fmt.Printf("Part 1: %d\n", s)
}

func printSects(start *section) {
	for start != nil {
		for range start.count {
			if start.id == -1 {
				fmt.Printf(".")
			} else {
				fmt.Printf("%d", start.id)
			}
		}
		start = start.next
	}
	fmt.Printf("\n")
}

func part2(input string) {
	sects := make(map[int]*section, len(input))
	var prev *section
	files := []*section{}
	for i, c := range input {
		if c == '\n' {
			continue
		}
		num, _ := strconv.Atoi(string(c))
		if i%2 == 0 {
			s := &section{
				id:    (i / 2),
				count: num,
				prev:  prev,
			}
			if s.prev != nil {
				s.prev.next = s
			}
			files = append(files, s)
			prev = s
			sects[s.id] = s
		} else {
			s := &section{
				id:    -1,
				count: num,
				prev:  prev,
			}
			if s.prev != nil {
				s.prev.next = s
			}
			prev = s
			sects[s.id] = s
		}
	}
	slices.Reverse(files)
	for _, f := range files {
		for empty := sects[0]; empty != nil && empty.id != f.id; empty = empty.next {
			if empty.id != -1 {
				continue
			}
			if empty.count > f.count {
				newEmptyReplacement := &section{
					id:    -1,
					count: f.count,
					next:  f.next,
					prev:  f.prev,
				}
				oldEmptyPrev := empty.prev
				empty.count -= f.count
				empty.prev = f
				oldFPrev := f.prev
				oldFPrev.next = newEmptyReplacement
				f.prev = oldEmptyPrev
				oldFNext := f.next
				f.next = empty
				if oldFNext != nil {
					oldFNext.prev = newEmptyReplacement
				}
				oldEmptyPrev.next = f
				break
			} else if empty.count == f.count {
				empty.id = f.id
				f.id = -1
				break
			}
		}
	}
	pos := 0
	s := 0
	curr := sects[0]
	for curr != nil {
		if curr.id == -1 {
			pos += curr.count
		} else {
			for range curr.count {
				s += int(curr.id) * pos
				pos++
			}
		}
		curr = curr.next
	}
	fmt.Printf("Part 2: %d\n", s)
}
