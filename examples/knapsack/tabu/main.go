package main

import (
	"fmt"

	"github.com/ccssmnn/hego"
	"github.com/ccssmnn/hego/mutate"
)

func knapsack(selection []bool, values, weights []float64, maxWeight float64) float64 {
	value := 0.0
	weight := 0.0
	for index, choice := range selection {
		if choice {
			value += values[index]
			weight += weights[index]
		}
	}
	if weight > maxWeight {
		value -= 100 * (weight - maxWeight)
	}
	return value
}

var values = []float64{69, 21, 33, 50, 89, 58, 27, 95, 52, 68, 26, 84, 46, 28, 25, 81, 82, 27, 50, 61}
var weights = []float64{6, 1, 1, 4, 9, 7, 3, 5, 7, 7, 9, 4, 4, 4, 8, 7, 7, 6, 5, 3}
var maxWeight = 30.0

// state is a selection boolean slice. It will implement the State interface for tabu search
type state []bool

// Neighbor produces a neighbor of the current state by flipping one selection
func (k state) Neighbor() hego.TabuState {
	return state(mutate.Flip(k))
}

// Equal returns true if two states are equal
func (k state) Equal(other hego.TabuState) bool {
	otherState := other.(state)
	for i := range k {
		if k[i] != otherState[i] {
			return false
		}
	}
	return true
}

// Objective returns the current objective from the knapsack problem. Negative, because lower is better
func (k state) Objective() float64 {
	return -knapsack(k, values, weights, maxWeight)
}

func main() {
	// select initial state
	initialState := state{false, true, true, false, false, true, false, false, false, true, true, false, false, true, false, true, false, true, false, false}

	// set algorithm parameters here
	settings := hego.TSSettings{}
	settings.MaxIterations = 100
	settings.Verbose = settings.MaxIterations / 10
	settings.TabuListSize = 50
	settings.NeighborhoodSize = 30

	// start tabu search main algorithm
	res, err := hego.TS(initialState, settings)

	if err != nil {
		fmt.Printf("Got error while running Anneal: %v", err)
		return
	}
	fmt.Printf("The solution found has an objective of %v \n", res.BestObjective)
}
