package main

import (
	"fmt"
	"math"

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
	// the penalty term is missing here since we dont allow elements that dont
	// fit to be selected
	return value
}

var values = []float64{69, 21, 33, 50, 89, 58, 27, 95, 52, 68, 26, 84, 46, 28, 25, 81, 82, 27, 50, 61}
var weights = []float64{6, 1, 1, 4, 9, 7, 3, 5, 7, 7, 9, 4, 4, 4, 8, 7, 7, 6, 5, 3}
var maxWeight = 30.0
var pheromones []float64
var bestPerformance = math.MaxFloat64

type ant struct {
	weight    float64
	value     float64
	selection []bool
}

// Init resets weight, value and selection. Is called before an ant is searching for another selection
func (a *ant) Init() {
	a.weight = 0.0
	a.value = 0.0
	a.selection = make([]bool, len(weights))
}

// Step adds next item to the selection. Here we also compute weight and value. Done is returned, when weight limit
// will be reached with another selection or no items are left to choose
func (a *ant) Step(next int) bool {
	a.weight += weights[next]
	a.value += values[next]
	a.selection[next] = true
	done := true
	for i, choice := range a.selection {
		if !choice {
			if a.weight+weights[i] < maxWeight {
				done = false
			}
		}
	}
	return done
}

// PerceivePheromone returns nonzero values for each unselected item and for any item that would not fit into the bag
func (a *ant) PerceivePheromone() []float64 {
	res := make([]float64, len(pheromones))
	copy(res, pheromones)
	// do not take items that are already taken
	for i, choice := range a.selection {
		if choice {
			res[i] = 0.0
		}
	}
	// do not take items, if their weight would increase load too much
	for i := range pheromones {
		if a.weight+weights[i] > maxWeight {
			res[i] = 0.0
		}
	}
	return res
}

// DropPheromone increases pheromone amount by 0.2 not considering the performance
func (a *ant) DropPheromone(performance float64) {
	for index, choice := range a.selection {
		if choice {
			pheromones[index] += 1 / (1 + (bestPerformance-performance)/bestPerformance)
		}
	}
}

// Evaporate applies factor and min to pheromone vector
func (a *ant) Evaporate(factor, min float64) {
	for i := range pheromones {
		pheromones[i] = math.Max(min, pheromones[i]*factor)
	}
}

// performance returns negative knapsack score
func (a *ant) Performance() float64 {
	performance := -knapsack(a.selection, values, weights, maxWeight)
	if performance < bestPerformance {
		bestPerformance = performance
	}
	return performance
}

func main() {
	initialPheromone := 1.0
	pheromones = make([]float64, len(weights))
	for i := range pheromones {
		pheromones[i] = initialPheromone
	}
	population := make([]hego.Ant, 100)
	for i := range population {
		population[i] = &ant{}
	}

	settings := hego.ACOSettings{}
	settings.Evaporation = 0.99
	settings.MinPheromone = 0.01
	settings.MaxIterations = 10
	settings.Verbose = settings.MaxIterations / 10 // log 10 steps to look at convergence behaviour

	result, err := hego.ACO(population, settings)

	if err != nil {
		fmt.Printf("Got error while running Ant Colony Optimization: %v", err)
	} else {
		fmt.Printf("Finished Ant Colony Optimization in %v! Needed %v function evaluations\n", result.Runtime, result.FuncEvaluations)
	}
	return
}
