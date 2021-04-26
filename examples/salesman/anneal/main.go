package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"strconv"
	"strings"

	"github.com/ccssmnn/hego"
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

// state represents a tour of cities
type state []int

// Neighbor produces a similar tour by swapping two cities in a tour
func (s state) Neighbor() hego.AnnealingState {
	return state(mutate.Swap(s))
}

// Energy counts the total tour length
func (s state) Energy() float64 {
	cost := 0.0
	position := s[0]
	for _, next := range s {
		cost += distances[position][next]
		position = next
	}
	cost += distances[position][s[0]]
	return cost
}

func main() {
	// read problem file
	err := readDistances()
	if err != nil {
		fmt.Printf("failed to read distances: %v", err)
		return
	}

	// produce one initial tour
	initialState := make(state, 48)
	for i := range initialState {
		initialState[i] = i
	}
	rand.Shuffle(len(initialState), func(i, j int) {
		initialState[i], initialState[j] = initialState[j], initialState[i]
	})

	// set algorithm parameters
	settings := hego.SASettings{}
	settings.MaxIterations = 1000000
	settings.Verbose = settings.MaxIterations / 10 // log 10 times during the process
	settings.Temperature = 10000.0                 // choose temperature in the range of initial random guesses energies
	settings.AnnealingFactor = 0.99999             // choose AnnealingFactor to make temperature reach low values at the end of the process

	// start simulated annealing
	result, err := hego.SA(&initialState, settings)

	if err != nil {
		fmt.Printf("Got error while running Anneal: %v", err)
	}
	finalEnergy := result.Energy
	fmt.Printf("Finished Simulated Annealing in %v! Tour Length: %v \n", result.Runtime, finalEnergy)
}
