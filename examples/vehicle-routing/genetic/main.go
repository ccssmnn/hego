package main

import (
	"fmt"
	"math/rand"

	"github.com/ccssmnn/hego"
	"github.com/ccssmnn/hego/crossover"
)

var nCities int
var depot int
var nVehicles int

var distances = [17][17]float64{
	{0, 548, 776, 696, 582, 274, 502, 194, 308, 194, 536, 502, 388, 354, 468, 776, 662},
	{548, 0, 684, 308, 194, 502, 730, 354, 696, 742, 1084, 594, 480, 674, 1016, 868, 1210},
	{776, 684, 0, 992, 878, 502, 274, 810, 468, 742, 400, 1278, 1164, 1130, 788, 1552, 754},
	{696, 308, 992, 0, 114, 650, 878, 502, 844, 890, 1232, 514, 628, 822, 1164, 560, 1358},
	{582, 194, 878, 114, 0, 536, 764, 388, 730, 776, 1118, 400, 514, 708, 1050, 674, 1244},
	{274, 502, 502, 650, 536, 0, 228, 308, 194, 240, 582, 776, 662, 628, 514, 1050, 708},
	{502, 730, 274, 878, 764, 228, 0, 536, 194, 468, 354, 1004, 890, 856, 514, 1278, 480},
	{194, 354, 810, 502, 388, 308, 536, 0, 342, 388, 730, 468, 354, 320, 662, 742, 856},
	{308, 696, 468, 844, 730, 194, 194, 342, 0, 274, 388, 810, 696, 662, 320, 1084, 514},
	{194, 742, 742, 890, 776, 240, 468, 388, 274, 0, 342, 536, 422, 388, 274, 810, 468},
	{536, 1084, 400, 1232, 1118, 582, 354, 730, 388, 342, 0, 878, 764, 730, 388, 1152, 354},
	{502, 594, 1278, 514, 400, 776, 1004, 468, 810, 536, 878, 0, 114, 308, 650, 274, 844},
	{388, 480, 1164, 628, 514, 662, 890, 354, 696, 422, 764, 114, 0, 194, 536, 388, 730},
	{354, 674, 1130, 822, 708, 628, 856, 320, 662, 388, 730, 308, 194, 0, 342, 422, 536},
	{468, 1016, 788, 1164, 1050, 514, 514, 662, 320, 274, 388, 650, 536, 342, 0, 764, 194},
	{776, 868, 1552, 560, 674, 1050, 1278, 742, 1084, 810, 1152, 274, 388, 422, 764, 0, 798},
	{662, 1210, 754, 1358, 1244, 708, 480, 856, 514, 468, 354, 844, 730, 536, 194, 798, 0},
}

// genome has two lists. The order of cities, and the assignment of a vehicle to a city.
// A tour is all cities assigned to a vehicle in the order they appear in the order list.
type genome struct {
	order      []int
	assignment []int
}

// assembleTour creates a tour list for the given vehicle index
func (g genome) assembleTour(vehicle int) []int {
	tour := make([]int, 0, nCities)
	for i, city := range g.order {
		if g.assignment[i] == vehicle {
			tour = append(tour, city)
		}
	}
	return tour
}

// Crossover combines the assignments and the order of both genomes
func (g genome) Crossover(other hego.Genome) hego.Genome {
	res := genome{
		assignment: make([]int, nCities),
		order:      make([]int, nCities),
	}
	// check which vehicle is assigned to which city
	assignments := make([][2]int, nCities)
	for i := range res.assignment {
		assignments[i] = [2]int{-1, -1}
	}
	for i := range g.assignment {
		gCity := g.order[i]
		assignments[gCity][0] = g.assignment[i]
		oCity := other.(genome).order[i]
		assignments[oCity][1] = other.(genome).assignment[i]
	}

	// combine city order
	res.order = crossover.OnePointPerm(g.order, other.(genome).order)

	// randomly select vehicle for a city depending on genomes
	for i, city := range res.order {
		if rand.Float64() > 0.5 {
			res.assignment[i] = assignments[city][0]
		} else {
			res.assignment[i] = assignments[city][1]
		}
	}
	return other
}

// Mutate changes the assignment of a city or the order of a city in a tour
func (g genome) Mutate() hego.Genome {
	if rand.Float64() > 0.5 {
		// change assignment
		i, j := rand.Intn(nCities), rand.Intn(nVehicles)
		if g.assignment[i] != j {
			g.assignment[i] = j
		} else {
			g.assignment[i] = (j + 1) % nVehicles
		}
	} else {
		// swap cities
		for {
			vehicle := rand.Intn(nVehicles)
			indices := make([]int, 0, nCities)
			for i, assignment := range g.assignment {
				if assignment == vehicle {
					indices = append(indices, i)
				}
			}
			if len(indices) <= 1 {
				continue
			}

			i, j := rand.Intn(len(indices)), rand.Intn(len(indices))
			if i == j {
				j = (i + 1) % len(indices)
			}
			ii, jj := indices[i], indices[j]
			g.order[ii], g.order[jj] = g.order[jj], g.order[ii]
			break
		}
	}
	return g
}

// Fitness is the total tour length
func (g genome) Fitness() float64 {
	totalLength := 0.0
	for vehicle := 0; vehicle < nVehicles; vehicle++ {
		tour := g.assembleTour(vehicle)
		// compute length of tour for this vehicle
		position := depot
		for _, next := range tour {
			totalLength += distances[position][next]
			position = next
		}
		totalLength += distances[position][depot]
	}
	return totalLength
}

func main() {
	nCities = 17
	depot = 0
	nVehicles = 4
	initialOrder := make([]int, nCities)
	initialAssignment := make([]int, nCities)
	for i := range initialOrder {
		initialOrder[i] = i
		initialAssignment[i] = rand.Intn(nVehicles)
	}

	// set algorithm parameters
	settings := hego.GASettings{}
	settings.Selection = hego.RankBasedSelection
	settings.MutationRate = 0.5
	settings.Elitism = 5
	settings.MaxIterations = 1000
	settings.Verbose = settings.MaxIterations / 10
	populationSize := 200

	population := make([]hego.Genome, populationSize)
	for i := range population {
		individuum := genome{
			order: make([]int, nCities), assignment: make([]int, nCities)}
		copy(individuum.order, initialOrder)
		copy(individuum.assignment, initialAssignment)
		rand.Shuffle(nCities, func(i, j int) {
			individuum.order[i], individuum.order[j] = individuum.order[j], individuum.order[i]
		})
		rand.Shuffle(nCities, func(i, j int) {
			individuum.assignment[i], individuum.assignment[j] = individuum.assignment[j], individuum.assignment[i]
		})
		population[i] = individuum
	}

	// start Genetic Algorithm
	result, err := hego.GA(population, settings)

	if err != nil {
		fmt.Printf("Got error while running Genetic Algorithm: %v", err)
	}

	fmt.Printf("Finished Genetic Algorithm in %v! Tour Length: %v \n", result.Runtime, result.BestFitness)
}
