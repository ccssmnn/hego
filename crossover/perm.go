package crossover

import (
	"math/rand"
)

// OnePointPerm cuts a in two pieces and fills the gap with values from b
// while preserving order. 12345678 + 26371485 -> 1234**** + *6*7**85 -> 12346785
func OnePointPerm(a, b []int) []int {
	if len(a) != len(b) {
		panic("expected inputs to have same length")
	}
	c := make([]int, len(a))
	cut := rand.Intn(len(c))
	// take every value before cut from a
	taken := map[int]bool{}
	for i := 0; i < cut; i++ {
		c[i] = a[i]
		taken[a[i]] = true
	}
	// return index of next untaken value in b
	nextFromB := func() int {
		for bindex := 0; bindex < len(b); bindex++ {
			_, exists := taken[b[bindex]]
			if !exists {
				return bindex
			}
		}
		panic("No untaken values in b left but another value was requested. Verify that the inputs have unique contents.")
	}
	// fill gaps in c with untaken values from b
	for i := cut; i < len(c); i++ {
		nextBIndex := nextFromB()
		taken[b[nextBIndex]] = true
		c[i] = b[nextBIndex]
	}
	return c
}

// TwoPointPerm takes a slice of a and fills the gaps with values from b
// while preserving order. 12345678 + 26371485 -> **3456** + 2**71*8* -> 27345618
func TwoPointPerm(a, b []int) []int {
	if len(a) != len(b) {
		panic("expected inputs to have same length")
	}
	c := make([]int, len(a))
	start, end := rand.Intn(len(c)), rand.Intn(len(c))
	if start > end {
		start, end = end, start
	}
	// take every value between start and end from a
	taken := map[int]bool{}
	for i := range c {
		if start <= i && i < end {
			c[i] = a[i]
			taken[a[i]] = true
		}
	}
	// return index of next untaken value in b
	nextFromB := func() int {
		for bindex := 0; bindex < len(b); bindex++ {
			_, exists := taken[b[bindex]]
			if !exists {
				return bindex
			}
		}
		panic("No untaken values in b left but another value was requested. Verify that the inputs have unique contents.")
	}
	// fill gaps in c with untaken values from b
	for i := range c {
		if i < start || end <= i {
			nextBIndex := nextFromB()
			taken[b[nextBIndex]] = true
			c[i] = b[nextBIndex]
		}
	}
	return c
}

// OnePointInt randomly chooses intersection point. Takes all elements from
// a, where the index is below the intersection point and b for the rest
// panics if length is different
func OnePointInt(a, b []int) []int {
	if len(a) != len(b) {
		panic("expected slices to have same length")
	}
	child := make([]int, len(a))
	copy(child, a)
	for i := rand.Intn(len(a)); i < len(a); i++ {
		child[i] = b[i]
	}
	return child
}

// TwoPointInt is analogue to OnePointInt with two intersection points
func TwoPointInt(a, b []int) []int {
	if len(a) != len(b) {
		panic("expected slices to have same length")
	}
	child := make([]int, len(a))
	copy(child, a)
	start, end := rand.Intn(len(a)), rand.Intn(len(a))
	if start > end {
		start, end = end, start
	}
	for i := start; i < end; i++ {
		child[i] = b[i]
	}
	return child
}
