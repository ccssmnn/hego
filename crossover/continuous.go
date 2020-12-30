package crossover

import (
	"math/rand"
)

// ArithmeticCrossover uniformly chooses a value between a and b elementwise
func ArithmeticCrossover(a, b []float64) []float64 {
	c := make([]float64, len(a))
	if len(a) != len(b) {
		panic("input vectors do not have the same length")
	}
	for i := range a {
		u := rand.Float64()
		c[i] = a[i] + (b[i]-a[i])*u
	}
	return c
}
