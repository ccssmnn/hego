package mutate

import "testing"

func TestGauss(t *testing.T) {
	a := []float64{1.0, 2.0, 3.0, 4.0}
	Gauss(a, 1.0)
}
