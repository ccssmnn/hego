package hego

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
	position2 := position1 + 1
	if position2 == len(x) {
		position2 = 0
	}
	neighbor := make([]int, len(x))
	copy(neighbor, x)
	neighbor[position1], neighbor[position2] = x[position2], x[position1]
	return neighbor
}

func contains(positions []int, position int) bool {
	for _, pos := range positions {
		if pos == position {
			return true
		}
	}
	return false
}

// Flip produces a neighbor of state by changing the value of one bit
func Flip(state []bool) []bool {
	position := rand.Intn(len(state))
	neighbor := make([]bool, len(state))
	copy(neighbor, state)
	neighbor[position] = !state[position]
	return neighbor
}

// Flipn produces a neighbor of state by changing the value of n bits
func Flipn(state []bool, n int) []bool {
	positions := []int{}
	for i := 0; i < n && i < len(state); i++ {
		nextPosition := rand.Intn(len(state))
		for contains(positions, nextPosition) {
			nextPosition = rand.Intn(len(state))
		}
		positions = append(positions, nextPosition)
	}

	neighbor := make([]bool, len(state))
	copy(neighbor, state)

	for _, position := range positions {
		neighbor[position] = !state[position]
	}

	return neighbor
}
