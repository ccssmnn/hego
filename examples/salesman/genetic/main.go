package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"strconv"
	"strings"

	"github.com/ccssmnn/hego"
	"github.com/ccssmnn/hego/crossover"
	"github.com/ccssmnn/hego/mutate"
)

var distances = [48][48]float64{}

func readDistances() error {
	file, err := ioutil.ReadFile("../att48.txt")
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	lines := strings.Split(string(file), "\n")
	if len(lines) != 48 {
		return fmt.Errorf("file has wrong number of lines. Wanted 48, got %v", len(lines))
	}
	for row, line := range lines {
		elems := strings.Split(line, " ")
		col := 0
		for _, elem := range elems {
			if len(elem) > 0 {
				distance, _ := strconv.Atoi(elem)
				distances[row][col] = float64(distance)
				col++
			}
		}
	}
	return nil
}

type genome struct {
	tour []int
}

func (g *genome) GetGene() []interface{} {
	gene := make([]interface{}, len(g.tour))
	for i := range gene {
		gene[i] = g.tour[i]
	}
	return gene
}

func (g *genome) Mutate() hego.GeneticGenome {
	new := &genome{}
	new.tour = mutate.Swap(g.tour)
	return new
}

func (g *genome) Crossover(other hego.GeneticGenome) hego.GeneticGenome {
	new := &genome{}
	gene := hego.ConvertInt(other.GetGene())
	new.tour = crossover.TwoPointPerm(g.tour, gene)
	return new
}

func (g *genome) Fitness() float64 {
	cost := 0.0
	position := g.tour[0]
	for _, next := range g.tour {
		cost += distances[position][next]
		position = next
	}
	cost += distances[position][g.tour[0]]
	return cost
}

func main() {
	err := readDistances()
	if err != nil {
		fmt.Printf("failed to read distances: %v", err)
		return
	}
	tour := make([]int, 0)
	for i := 0; i < 48; i++ {
		tour = append(tour, i)
	}

	population := make([]hego.GeneticGenome, 100)
	for i := range population {
		t := make([]int, len(tour))
		copy(t, tour)
		rand.Shuffle(len(t), func(i, j int) { t[i], t[j] = t[j], t[i] })
		population[i] = &genome{tour: t}
	}

	settings := hego.GeneticSettings{}
	settings.Selection = hego.RankBasedSelection
	settings.MutationRate = 0.5
	settings.Elitism = 1
	settings.MaxIterations = 10000
	settings.Verbose = 1000

	result, err := hego.Genetic(population, settings)

	if err != nil {
		fmt.Printf("Got error while running Genetic Algorithm: %v", err)
	} else {
		fmt.Printf("Finished Genetic Algorithm in %v! Needed %v function evaluations\n", result.Runtime, result.FuncEvaluations)
	}
	return
}
