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

// state is a two element vector, it will implement the State interface for Simulated Annealing
type state []float64

// Neighbor produces another state by adding gaussian noise to the current state
func (s state) Neighbor() hego.AnnealingState {
	return state(mutate.Gauss(s, 0.3))
}

// Energy returns the energy of the current state. Lower is better
func (s state) Energy() float64 {
	return rastringin(s[0], s[1])
}

func main() {

	initialState := state{5.0, 5.0}

	settings := hego.SASettings{}
	settings.MaxIterations = 100000
	settings.Verbose = 10000
	settings.Temperature = 10.0       // choose temperature in the range of the systems energy
	settings.AnnealingFactor = 0.9999 // decrementing the temperature leads to convergence, we want to reach convergence when approaching the end of iterations

	// start simulated annealing algorithm
	result, err := hego.SA(initialState, settings)

	if err != nil {
		fmt.Printf("Got error while running Anneal: %v", err)
	}
	fmt.Printf("Finished Simulated Annealing in %v! Result: %v, Value: %v \n", result.Runtime, result.States[result.Iterations], result.Energies[result.Iterations])
	return
}
