package mutate

import "math/rand"

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
