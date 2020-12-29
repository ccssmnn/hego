package mutate

import "math/rand"

// Swap produces a neighbor of x by swapping two values
func Swap(x []int) []int {
	position1 := rand.Intn(len(x))
	position2 := rand.Intn(len(x))
	neighbor := make([]int, len(x))
	copy(neighbor, x)
	neighbor[position1], neighbor[position2] = x[position2], x[position1]
	return neighbor
}

// SwapClose produces a neighbor of x by swapping two values that are next to each other
func SwapClose(x []int) []int {
	position1 := rand.Intn(len(x))
	position2 := (position1 + 1) % len(x)
	neighbor := make([]int, len(x))
	copy(neighbor, x)
	neighbor[position1], neighbor[position2] = x[position2], x[position1]
	return neighbor
}
