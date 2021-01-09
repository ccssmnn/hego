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

// state is a selection boolean slice. It will implement the State interface for simulated annealing
type state []bool

// Neighbor produces a neighbor of the current state by flipping one selection
func (k state) Neighbor() hego.AnnealingState {
	return state(mutate.Flip(k))
}

// Energy returns the current objective from the knapsack problem. Negative, because lower is better
func (k state) Energy() float64 {
	return -knapsack(k, values, weights, maxWeight)
}

func main() {
	// select initial state
	initialState := state{false, true, true, false, false, true, false, false, false, true, true, false, false, true, false, true, false, true, false, false}

	// set algorithm parameters here. Temperature and AnnealingFactor are critical for
	// the convergence behaviour
	settings := hego.SASettings{}
	settings.MaxIterations = 1000
	settings.Verbose = 100 // log status every 100 steps. very useful for choosing parameters
	settings.Temperature = 1000.0
	settings.AnnealingFactor = 0.99

	// start simulated annealing main algorithm
	res, err := hego.SA(&initialState, settings)

	if err != nil {
		fmt.Printf("Got error while running Anneal: %v", err)
		return
	}
	// extract result
	solution := res.States[len(res.States)-1].Energy()
	fmt.Printf("The solution found has an energy of %v \n", solution)
}
