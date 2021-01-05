package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"

	"github.com/ccssmnn/hego"
)

var distances = [][]float64{}

func readDistances() error {
	file, err := ioutil.ReadFile("../att48.txt")
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	lines := strings.Split(string(file), "\n")
	if len(lines) != 48 {
		return fmt.Errorf("file has wrong number of lines. Wanted 48, got %v", len(lines))
	}
	for _, line := range lines {
		dLine := make([]float64, len(lines))
		elems := strings.Split(line, " ")
		col := 0
		for _, elem := range elems {
			if len(elem) > 0 {
				distance, _ := strconv.Atoi(elem)
				dLine[col] = float64(distance)
				col++
			}
		}
		distances = append(distances, dLine)
	}
	return nil
}

var pheromones = [][]float64{}

type ant struct {
	position int
	tour     []int
}

func (a *ant) Init() {
	// always start at city 0
	a.position = 0
	// reset tour
	a.tour = []int{0}
}

func (a *ant) Step(next int) bool {
	a.position = next
	a.tour = append(a.tour, next)
	return len(a.tour) == 48
}

func (a *ant) PerceivePheromone() []float64 {
	p := pheromones[a.position]
	res := make([]float64, len(p))
	copy(res, p)
	for _, stop := range a.tour {
		res[stop] = 0.0
	}
	res[a.position] = 0.0
	return res
}

func (a *ant) DropPheromone(performance float64) {
	for i := range a.tour {
		if i == len(a.tour)-1 {
			continue
		}
		prev, next := a.tour[i], a.tour[i+1]
		pheromones[prev][next] += 0.2
		pheromones[next][prev] += 0.2
	}
}

func (a *ant) Evaporate(factor, min float64) {
	for i := range pheromones {
		for j := range pheromones[i] {
			pheromones[i][j] = math.Max(min, pheromones[i][j]*factor)
		}
	}
}

func (a *ant) Performance() float64 {
	cost := 0.0
	position := a.tour[0]
	for _, next := range a.tour {
		cost += distances[position][next]
		position = next
	}
	cost += distances[position][a.tour[0]]
	return cost
}

func main() {
	err := readDistances()
	if err != nil {
		fmt.Printf("failed to read distances: %v", err)
		return
	}
	pheromones = hego.InitializePheromoneMatrix(len(distances), 1.0)
	tour := make([]int, 0)
	for i := 0; i < 48; i++ {
		tour = append(tour, i)
	}

	population := make([]hego.Ant, 100)
	for i := range population {
		population[i] = &ant{tour: []int{}, position: -1}
	}

	settings := hego.ACOSettings{}
	settings.Evaporation = 0.9
	settings.MinPheromone = 0.001
	settings.MaxIterations = 1000
	settings.Verbose = settings.MaxIterations / 10

	result, err := hego.ACO(population, settings)

	if err != nil {
		fmt.Printf("Got error while running Genetic Algorithm: %v", err)
	} else {
		fmt.Printf("Finished Genetic Algorithm in %v! Needed %v function evaluations\n", result.Runtime, result.FuncEvaluations)
	}
	return
}
