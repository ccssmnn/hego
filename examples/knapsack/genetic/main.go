package main

import (
	"fmt"
	"math/rand"

	"github.com/ccssmnn/hego"
	"github.com/ccssmnn/hego/crossover"
	"github.com/ccssmnn/hego/mutate"
)

// knapsack evaluates the objective for a given selection. Adds penalties for exceeding
// maxWeight
func knapsack(selection []bool, values, weights []float64, maxWeight float64) float64 {
	value := 0.0
	weight := 0.0
	for index, choice := range selection {
		if choice {
			value += values[index]
			weight += weights[index]
		}
	}
	// penalty
	if weight > maxWeight {
		value -= 100 * (weight - maxWeight)
	}
	return value
}

// parameters of current knapsack problem
var values = []float64{69, 21, 33, 50, 89, 58, 27, 95, 52, 68, 26, 84, 46, 28, 25, 81, 82, 27, 50, 61}
var weights = []float64{6, 1, 1, 4, 9, 7, 3, 5, 7, 7, 9, 4, 4, 4, 8, 7, 7, 6, 5, 3}
var maxWeight = 30.0

// genome is a vector of bool. True represents a selected item
type genome []bool

// Crossover uses uniform Crossover for booleans
func (k genome) Crossover(other hego.Genome) hego.Genome {
	return genome(crossover.UniformBool(k, other.(genome)))
}

// Mutate flips one bit, e.g. makes one selected item unselected or vice versa
func (k genome) Mutate() hego.Genome {
	return genome(mutate.Flip(k))
}

// Fitness returns the negative knapsack score. Lower is better
func (k genome) Fitness() float64 {
	return -knapsack(k, values, weights, maxWeight)
}

func main() {
	// initialize population
	populationSize := 100
	population := make([]hego.Genome, populationSize)
	for i := range population {
		individuum := make(genome, len(values))
		for j := range values {
			individuum[j] = rand.Int31()%2 == 0 // efficient random boolean
		}
		population[i] = individuum
	}

	// set algorithm parameterrs
	settings := hego.GASettings{}
	settings.MutationRate = 0.3
	settings.Elitism = 1
	settings.MaxIterations = 100
	settings.Verbose = 10

	result, err := hego.GA(population, settings)
	if err != nil {
		fmt.Printf("Got error while running Anneal: %v", err)
		return
	}
	fmt.Printf("Finished Genetic Algorithm in %v! Needed %v function evaluations\n", result.Runtime, result.FuncEvaluations)
	return
}
