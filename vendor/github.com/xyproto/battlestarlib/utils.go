package battlestarlib

import (
	"math"
	"strconv"
)

// Find the position of a string in a list of strings, -1 if not found
func pos(sl []string, s string) int {
	for i, e := range sl {
		if e == s {
			return i
		}
	}
	return -1
}

// Maps the function f over a slice of strings
func maps(sl []string, f func(string) string) []string {
	newl := make([]string, len(sl))
	for i, element := range sl {
		newl[i] = f(element)
	}
	return newl
}

// Checks if a slice of strings has the given string
func has(sl []string, s string) bool {
	for _, e := range sl {
		if e == s {
			return true
		}
	}
	return false
}

// Checks if a slice of ints has the given int
func hasi(il []int, i int) bool {
	for _, e := range il {
		if e == i {
			return true
		}
	}
	return false
}

// Given a non-hex number as a string, like "123", return the number of bits of space it takes.
// For the case of "123" the answer would be 7.
// Return 0 if it's not a number
func numbits(number string) int {
	n, err := strconv.Atoi(number)
	if err != nil {
		return 0
	}
	return int(math.Ceil(math.Log2(float64(n))))
}
