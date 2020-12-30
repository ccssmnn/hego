package crossover

import "testing"

func TestArithmeticCrossover(t *testing.T) {
	a := []float64{1.0, 1.0}
	b := []float64{2.0, 2.0}
	c := ArithmeticCrossover(a, b)
	for i, value := range c {
		if (value-a[i])/(b[i]-a[i]) > 1.0 {
			t.Error("unexpected range in created value")
		}
	}
}
