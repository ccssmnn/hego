package crossover

import "testing"

// findInSlice counts appearances of value in slice
func findInSlice(value int, slice []int) int {
	count := 0
	for _, v := range slice {
		if v == value {
			count++
		}
	}
	return count
}

func TestTwoPointPerm(t *testing.T) {
	a := []int{1, 2, 3, 4, 5, 6, 7, 8}
	b := []int{8, 7, 6, 5, 4, 3, 2, 1}
	c := TwoPointPerm(a, b)
	for _, v := range a {
		appearances := findInSlice(v, c)
		if appearances != 1 {
			t.Errorf("unexpected number of appearances in crossover result. Got %v, expected 1", appearances)
		}
	}
	d := []int{1, 2, 3}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("crossover method should panic for inputs with unequal lengths")
		}
	}()
	TwoPointPerm(a, d)
}

func TestTwoPointPerm2(t *testing.T) {
	a := []int{1, 2, 3, 4, 5, 6, 7, 8}
	b := []int{1, 1, 1, 1, 1, 1, 1, 1}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("permutation based crossover should fail for duplicates in the input")
		}
	}()
	TwoPointPerm(a, b)
}

func TestOnePointPerm(t *testing.T) {
	a := []int{1, 2, 3, 4, 5, 6, 7, 8}
	b := []int{8, 7, 6, 5, 4, 3, 2, 1}
	c := OnePointPerm(a, b)
	for _, v := range a {
		appearances := findInSlice(v, c)
		if appearances != 1 {
			t.Errorf("unexpected number of appearances in crossover result. Got %v, expected 1", appearances)
		}
	}
	d := []int{1, 2, 3}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("crossover method should panic for inputs with unequal lengths")
		}
	}()
	OnePointPerm(a, d)
}

func TestOnePointPerm2(t *testing.T) {
	a := []int{1, 2, 3, 4, 5, 6, 7, 8}
	b := []int{1, 1, 1, 1, 1, 1, 1, 1}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("permutation based crossover should fail for duplicates in the input")
		}
	}()
	OnePointPerm(a, b)
}
