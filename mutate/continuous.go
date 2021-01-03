package mutate

import (
	"math/rand"
)

// Gauss returns a new vector with gaussian noise. dev is custom deviation
func Gauss(a []float64, dev float64) []float64 {
	res := make([]float64, len(a))
	copy(res, a)
	for i := range a {
		res[i] = a[i] + rand.NormFloat64()*dev
	}
	return res
}
