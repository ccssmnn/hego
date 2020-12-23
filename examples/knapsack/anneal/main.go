package main

import (
	"fmt"

	"github.com/ccssmnn/hego"
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
		value -= 10 * (weight - maxWeight)
	}
	return value
}

var values = []float64{69, 21, 33, 50, 89, 58, 27, 95, 52, 68, 26, 84, 46, 28, 25, 81, 82, 27, 50, 61}
var weights = []float64{6, 1, 1, 4, 9, 7, 3, 5, 7, 7, 9, 4, 4, 4, 8, 7, 7, 6, 5, 3}
var maxWeight = 30.0

type knapsackState struct {
	selection []bool
}

func (k *knapsackState) Clone() hego.State {
	clone := knapsackState{selection: make([]bool, len(k.selection))}
	copy(clone.selection, k.selection)
	return &clone
}

func (k *knapsackState) Neighbor() hego.State {
	n := knapsackState{}
	n.selection = hego.Flip(k.selection)
	return &n
}

func (k *knapsackState) Energy() float64 {
	return -knapsack(k.selection, values, weights, maxWeight)
}

func main() {
	initialState := knapsackState{
		selection: []bool{false, true, true, false, false, true, false, false, false, true, true, false, false, true, false, true, false, true, false, false},
	}

	settings := hego.Settings{
		MaxIterations: 100,
		Verbose:       10,
	}

	temperature := 100.0
	annealingFactor := 0.9

	result, err := hego.Anneal(&initialState, temperature, annealingFactor, settings)

	if err != nil {
		fmt.Printf("Got error while running Anneal: %v", err)
	} else {
		fmt.Printf("Finished Annealing in %v! Result: %v, Value: %v \n", result.Runtime, result.States[result.Iterations], -result.Energies[result.Iterations])
	}
	return
}
