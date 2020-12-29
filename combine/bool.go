package combine

import (
	"math/rand"
)

// UniformCrossoverBool uniformly selects attributes from a or b
// panics if a and b have different lengths
func UniformCrossoverBool(a, b []bool) []bool {
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

// OnePointCrossoverBool randomly chooses intersection point. Takes all elements from
// a, where the index is below the intersection point and b for the rest
// panics if length is different
func OnePointCrossoverBool(a, b []bool) []bool {
	if len(a) != len(b) {
		panic("expected slices to have same length")
	}
	child := make([]bool, len(a))
	copy(child, a)
	intersection := rand.Intn(len(a))
	for i := intersection; intersection < len(a); i++ {
		child[i] = b[i]
	}
	return child
}
