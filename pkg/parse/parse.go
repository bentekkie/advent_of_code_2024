package parse

import (
	"strconv"
	"strings"

	"golang.org/x/exp/constraints"
)

type num interface {
	constraints.Integer | constraints.Float
}

func MustAtoi[T num](s string) T {
	n, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return T(n)
}

func NumList[T num](s string, sep string) []T {
	var nums []T
	for _, part := range strings.Split(s, sep) {
		nums = append(nums, MustAtoi[T](part))
	}
	return nums
}
