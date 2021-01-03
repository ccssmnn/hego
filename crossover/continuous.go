package crossover

import (
	"math/rand"
)

// Arithmetic chooses u in range and performs c = u * a + (1-u) * b
// for every element in a and b
func Arithmetic(a, b []float64, uRange [2]float64) []float64 {
	c := make([]float64, len(a))
	if len(a) != len(b) {
		panic("input vectors do not have the same length")
	}
	for i := range a {
		u := uRange[0] + (uRange[1]-uRange[0])*rand.Float64()
		c[i] = a[i] + (b[i]-a[i])*u
	}
	return c
}
