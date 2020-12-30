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

func TestOrderBasedCrossover(t *testing.T) {
	a := []int{1, 2, 3, 4, 5, 6, 7, 8}
	b := []int{8, 7, 6, 5, 4, 3, 2, 1}
	c := OrderBasedCrossover(a, b)
	for _, v := range a {
		appearances := findInSlice(v, c)
		if appearances != 1 {
			t.Errorf("unexpected number of appearances in crossover result. Got %v, expected 1", appearances)
		}
	}
}
