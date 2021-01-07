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
var bestPerformance = math.MaxFloat64

type ant struct {
	length   float64
	position int
	tour     []int
}

func (a *ant) Init() {
	// always start at city 0
	a.position = 0
	a.length = 0.0
	// reset tour
	a.tour = []int{0}
}

// Step updates position, appends next stop to tour and returns done when tour length is equal to city count
func (a *ant) Step(next int) bool {
	a.length += distances[a.position][next]
	a.position = next
	a.tour = append(a.tour, next)
	done := len(a.tour) == 48
	if done {
		a.tour = append(a.tour, 0)
		a.length += distances[a.position][0]
		a.position = 0
	}
	return done
}

// PerceivePheromone returns values from pheromone matrix for current position
// and multiplies it by a 1/distance. This makes closer cities more attractive
// (nearest neighbor heuristic). Also all stops that have been visited already
// will have 0.0 pheromone
func (a *ant) PerceivePheromone() []float64 {
	p := pheromones[a.position]
	d := distances[a.position]
	res := make([]float64, len(p))
	copy(res, p)
	for i := range res {
		res[i] *= 1 / d[i] // nearest neighbor factor
	}
	for _, stop := range a.tour {
		res[stop] = 0.0 // set pheromone to 0 when city has been visited
	}
	return res
}

// DropPheromone adds pheromone to pheromone matrix along the tour of this ant
// the amount of pheromone that is dropped depends on the performance of the tour
// compared with the current best performance
func (a *ant) DropPheromone(performance float64) {
	for i := range a.tour {
		if i == len(a.tour)-1 {
			continue
		}
		prev, next := a.tour[i], a.tour[i+1]
		pheromones[prev][next] += 1 / (1 + (bestPerformance-performance)/bestPerformance)
		pheromones[next][prev] += 1 / (1 + (bestPerformance-performance)/bestPerformance)
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
	length := a.length
	if length < bestPerformance {
		bestPerformance = length
	}
	return length
}

func main() {
	err := readDistances()
	if err != nil {
		fmt.Printf("failed to read distances: %v", err)
		return
	}

	pheromones = make([][]float64, len(distances))
	initialPheromone := 1.0
	for i := range pheromones {
		line := make([]float64, len(distances[i]))
		for j := range line {
			if i != j {
				line[j] = initialPheromone
			}
		}
		pheromones[i] = line
	}
	tour := make([]int, 0)
	for i := 0; i < 48; i++ {
		tour = append(tour, i)
	}

	population := make([]hego.Ant, 100)
	for i := range population {
		population[i] = &ant{tour: []int{}, position: -1}
	}

	settings := hego.ACOSettings{}
	settings.Evaporation = 0.99
	settings.MinPheromone = 0.01
	settings.MaxIterations = 1000
	settings.Verbose = settings.MaxIterations / 10

	result, err := hego.ACO(population, settings)

	if err != nil {
		fmt.Printf("Got error while running Ant Colony Optimization: %v", err)
	} else {
		fmt.Printf("Finished Ant Colony Optimization in %v! Needed %v function evaluations\n", result.Runtime, result.FuncEvaluations)
	}
	return
}
