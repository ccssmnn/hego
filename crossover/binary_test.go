package crossover

import (
	"testing"
)

// only calls methods, since its difficult to test for randomized results
func TestUniformBool(t *testing.T) {
	a := []bool{true, true, true, true}
	b := []bool{false, false, false, false}
	UniformBool(a, b)
	d := []bool{true, true}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("crossover method should panic for inputs with unequal lengths")
		}
	}()
	UniformBool(a, d)
}

func TestOnePointBool(t *testing.T) {
	a := []bool{true, true, true, true}
	b := []bool{false, false, false, false}
	c := OnePointBool(a, b)
	bStart := -1
	for i := range c {
		if bStart == -1 {
			if c[i] == false {
				bStart = i
			}
		} else {
			if c[i] != false {
				t.Error("OnePointBool crossover returned unexpected result. All values after intersection point must come from b.")
			}
		}
	}
	d := []bool{true, true}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("crossover method should panic for inputs with unequal lengths")
		}
	}()
	OnePointBool(a, d)
}

func TestTwoPointBool(t *testing.T) {
	a := []bool{true, true, true, true}
	b := []bool{false, false, false, false}
	c := TwoPointBool(a, b)
	bStart := -1
	bEnd := -1
	for i := range c {
		if bStart == -1 && bEnd == -1 {
			if c[i] == false {
				bStart = i
			}
		} else if bStart != -1 && bEnd == -1 {
			if c[i] == true {
				bStart = i
			}
		} else {
			if c[i] != true {
				t.Error("TwoPointBool crossover returned unexpected result. All values after second intersection point must come from a.")
			}
		}
	}
	d := []bool{true, true}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("crossover method should panic for inputs with unequal lengths")
		}
	}()
	TwoPointBool(a, d)
}
