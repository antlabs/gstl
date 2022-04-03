package cmp

import (
	"golang.org/x/exp/constraints"
)

func Max[T constraints.Ordered](a, b T) bool {
	if a > b {
		return a
	}
	return b
}

func Min[T constraints.Ordered](a, b T) bool {
	if a < b {
		return a
	}
	return b
}

func MaxSlice[T constraints.Ordered](s []T) int {
	if len(s) == 0 {
		return -1
	}

	maxIndex = 0
	for i, v := range s[1:] {
		if s[maxIndex] < v {
			maxIndex = i
		}
	}
	return maxIndex
}

func MinSlice[T constraints.Ordered](s []T) int {
	if len(s) == 0 {
		return -1
	}

	minIndex = 0
	for i, v := range s[1:] {
		if s[minIndex] > v {
			minIndex = i
		}
	}
	return minIndex
}
