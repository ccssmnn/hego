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

// genome is a slice of integers. Each element encodes a city and the order encodes the tour
type genome []int

// Mutate uses Swap to swap two cities in the tour
func (g genome) Mutate() hego.Genome {
	return genome(mutate.Swap(g))
}

func (g genome) Crossover(other hego.Genome) hego.Genome {
	return genome(crossover.TwoPointPerm(g, other.(genome)))
}

func (g genome) Fitness() float64 {
	cost := 0.0
	position := g[0]
	for _, next := range g {
		cost += distances[position][next]
		position = next
	}
	cost += distances[position][g[0]]
	return cost
}

func main() {
	// read distances from text file
	err := readDistances()
	if err != nil {
		fmt.Printf("failed to read distances: %v", err)
		return
	}

	// initialTour is counting up from zero {0, 1, 2, 3, ...}
	initialTour := make([]int, 0)
	for i := 0; i < 48; i++ {
		initialTour = append(initialTour, i)
	}

	// these are the algorithm parameters to tweak for your problem
	settings := hego.GASettings{}
	settings.Selection = hego.RankBasedSelection
	settings.MutationRate = 0.5
	settings.Elitism = 5
	settings.MaxIterations = 1000
	settings.Verbose = 100
	populationSize := 200

	// initialize population
	population := make([]hego.Genome, populationSize)
	for i := range population {
		// create a new randomized tour and add it to population
		tour := make(genome, len(initialTour))
		copy(tour, initialTour)
		rand.Shuffle(len(tour), func(i, j int) { tour[i], tour[j] = tour[j], tour[i] })
		population[i] = tour
	}

	// perform genetic algorithm
	result, err := hego.GA(population, settings)

	if err != nil {
		fmt.Printf("Got error while running Genetic Algorithm: %v", err)
		return
	}
	fmt.Printf("Finished Genetic Algorithm in %v! Needed %v function evaluations\n", result.Runtime, result.FuncEvaluations)
}
