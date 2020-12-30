package main

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/ccssmnn/hego"
	"github.com/ccssmnn/hego/crossover"
)

func rastringin(x, y float64) float64 {
	return 10*2 + (x*x - 10*math.Cos(2*math.Pi*x)) + (y*y - 10*math.Cos(2*math.Pi*y))
}

type genome struct {
	v []float64
}

func (g *genome) GetGene() []interface{} {
	gene := make([]interface{}, len(g.v))
	for i, value := range g.v {
		gene[i] = value
	}
	return gene
}

func (g *genome) Crossover(other hego.GeneticGenome) hego.GeneticGenome {
	clone := genome{v: make([]float64, len(g.v))}
	gene := hego.ConvertFloat64(other.GetGene())
	clone.v = crossover.ArithmeticCrossover(g.v, gene)
	return &clone
}

func (g *genome) Mutate() hego.GeneticGenome {
	n := genome{v: make([]float64, len(g.v))}
	n.v[0] = g.v[0] + rand.NormFloat64()
	n.v[1] = g.v[1] + rand.NormFloat64()
	return &n
}

func (g *genome) Fitness() float64 {
	return rastringin(g.v[0], g.v[1])
}

func main() {
	population := make([]hego.GeneticGenome, 100)
	for i := range population {
		population[i] = &genome{v: []float64{-10.0 + 10.0*rand.Float64(), -10 + 10*rand.Float64()}}
	}

	settings := hego.GeneticSettings{}
	settings.MutationRate = 0.3
	settings.Elitism = 0
	settings.MaxIterations = 100
	settings.Verbose = 10

	result, err := hego.Genetic(population, settings)

	if err != nil {
		fmt.Printf("Got error while running Anneal: %v", err)
	} else {
		fmt.Printf("Finished Genetic Algorithm in %v! Needed %v function evaluations\n", result.Runtime, result.FuncEvaluations)
	}
	return
}
