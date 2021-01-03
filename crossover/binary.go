package crossover

import (
	"math/rand"
)

// UniformBool uniformly selects attributes from a or b
// panics if a and b have different lengths
func UniformBool(a, b []bool) []bool {
	if len(a) != len(b) {
		panic("expected slices to have same length")
	}
	child := make([]bool, len(a))
	for i := range a {
		if rand.Float64() > 0.5 {
			child[i] = a[i]
		} else {
			child[i] = b[i]
		}
	}
	return child
}

// OnePointBool randomly chooses intersection point. Takes all elements from
// a, where the index is below the intersection point and b for the rest
// panics if length is different
func OnePointBool(a, b []bool) []bool {
	if len(a) != len(b) {
		panic("expected slices to have same length")
	}
	child := make([]bool, len(a))
	copy(child, a)
	for i := rand.Intn(len(a)); i < len(a); i++ {
		child[i] = b[i]
	}
	return child
}

// TwoPointBool is analogue to OnePointBool with two intersection points
func TwoPointBool(a, b []bool) []bool {
	if len(a) != len(b) {
		panic("expected slices to have same length")
	}
	child := make([]bool, len(a))
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
