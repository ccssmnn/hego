package main

import (
	"fmt"
	"math"

	"github.com/ccssmnn/hego"
	"github.com/ccssmnn/hego/mutate"
)

func rastringin(x, y float64) float64 {
	return 10*2 + (x*x - 10*math.Cos(2*math.Pi*x)) + (y*y - 10*math.Cos(2*math.Pi*y))
}

// state is a two element vector, it will implement the State interface for Tabu Search
type state []float64

// Neighbor produces another state by adding gaussian noise to the current state
func (s state) Neighbor() hego.TabuState {
	return state(mutate.Gauss(s, 0.3))
}

// Equal returns true, of two states are almost equal
func (s state) Equal(other hego.TabuState) bool {
	otherState := other.(state)
	for i := range s {
		if math.Abs(s[i]-otherState[i]) > 1e-4 {
			return false
		}
	}
	return true
}

// Objective of the current state. Lower is better
func (s state) Objective() float64 {
	return rastringin(s[0], s[1])
}

func main() {

	initialState := state{5.0, 5.0}

	settings := hego.TSSettings{}
	settings.MaxIterations = 1000
	settings.Verbose = settings.MaxIterations / 10
	settings.TabuListSize = 50
	settings.NeighborhoodSize = 25

	// start tabu search algorithm
	result, err := hego.TS(initialState, settings)

	if err != nil {
		fmt.Printf("Got error while running Anneal: %v", err)
	}
	finalState := result.States[len(result.States)-1]
	finalEnergy := finalState.Objective()
	fmt.Printf("Finished Tabu Search in %v! Result: %v, Value: %v \n", result.Runtime, finalState, finalEnergy)
	return
}
