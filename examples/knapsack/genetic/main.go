package main

import (
	"fmt"
	"math/rand"

	"github.com/ccssmnn/hego"
	"github.com/ccssmnn/hego/crossover"
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

type knapsackGenome struct {
	selection []bool
}

func (k *knapsackGenome) GetGene() []interface{} {
	gene := make([]interface{}, len(k.selection))
	for i, value := range k.selection {
		gene[i] = value
	}
	return gene
}

func (k *knapsackGenome) Crossover(other hego.GeneticGenome) hego.GeneticGenome {
	new := knapsackGenome{selection: make([]bool, len(k.selection))}
	otherSelection := hego.ConvertBool(other.GetGene())
	new.selection = crossover.UniformBool(k.selection, otherSelection)
	return &new
}

func (k *knapsackGenome) Mutate() hego.GeneticGenome {
	n := knapsackGenome{}
	n.selection = mutate.Flip(k.selection)
	return &n
}

func (k *knapsackGenome) Fitness() float64 {
	return -knapsack(k.selection, values, weights, maxWeight)
}

func main() {
	population := make([]hego.GeneticGenome, 0)
	for i := 0; i < 10; i++ {
		individuum := knapsackGenome{selection: make([]bool, len(values))}
		for j := range values {
			individuum.selection[j] = rand.Float64() > 0.5
		}
		population = append(population, &individuum)
	}

	settings := hego.GeneticSettings{}
	settings.MutationRate = 0.1
	settings.Elitism = 0
	settings.MaxIterations = 10
	settings.Verbose = 10

	result, err := hego.Genetic(population, settings)

	if err != nil {
		fmt.Printf("Got error while running Anneal: %v", err)
	} else {
		fmt.Printf("Finished Genetic Algorithm in %v! Needed %v function evaluations\n", result.Runtime, result.FuncEvaluations)
	}
	return
}
