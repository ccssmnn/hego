package crossover

import "testing"

func TestArithmetic(t *testing.T) {
	a := []float64{1.0, 1.0}
	b := []float64{2.0, 2.0}
	c := Arithmetic(a, b, [2]float64{0.0, 1.0})
	for i, value := range c {
		if (value-a[i])/(b[i]-a[i]) > 1.0 {
			t.Error("unexpected range in created value")
		}
	}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("crossover method should panic for inputs with unequal lengths")
		}
	}()
	d := []float64{1.0, 2.0, 3.0}
	Arithmetic(a, d, [2]float64{0.0, 1.0})
}
