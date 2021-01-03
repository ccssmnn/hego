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

func (k *knapsackState) Neighbor() hego.AnnealState {
	n := knapsackState{}
	n.selection = mutate.Flip(k.selection)
	return &n
}

func (k *knapsackState) Energy() float64 {
	return -knapsack(k.selection, values, weights, maxWeight)
}

func main() {
	initialState := knapsackState{
		selection: []bool{false, true, true, false, false, true, false, false, false, true, true, false, false, true, false, true, false, true, false, false},
	}

	settings := hego.AnnealSettings{}
	settings.MaxIterations = 100
	settings.Verbose = 10
	settings.Temperature = 100.0
	settings.AnnealingFactor = 0.9

	result, err := hego.Anneal(&initialState, settings)

	if err != nil {
		fmt.Printf("Got error while running Anneal: %v", err)
	} else {
		fmt.Printf("Finished Simulated Annealing Algorithm in %v! Needed %v function evaluations\n", result.Runtime, result.FuncEvaluations)
	}
	return
}
