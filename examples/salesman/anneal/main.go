package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"strconv"
	"strings"
	"time"

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

type state struct {
	tour []int
}

func (s *state) Clone() hego.AnnealState {
	new := &state{tour: make([]int, len(s.tour))}
	copy(new.tour, s.tour)
	return new
}

func (s *state) Neighbor() hego.AnnealState {
	neighbor := &state{tour: []int{}}
	neighbor.tour = mutate.Swap(s.tour)
	return neighbor
}

func (s *state) Energy() float64 {
	cost := 0.0
	position := s.tour[0]
	for _, next := range s.tour {
		cost += distances[position][next]
		position = next
	}
	cost += distances[position][s.tour[0]]
	return cost
}

func main() {
	err := readDistances()
	if err != nil {
		fmt.Printf("failed to read distances: %v", err)
		return
	}
	initialState := state{
		tour: make([]int, 0),
	}
	for i := 0; i < 48; i++ {
		initialState.tour = append(initialState.tour, i)
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(initialState.tour), func(i, j int) {
		initialState.tour[i], initialState.tour[j] = initialState.tour[j], initialState.tour[i]
	})

	settings := hego.AnnealSettings{}
	settings.MaxIterations = 1000000
	settings.Verbose = 100000
	settings.Temperature = 100000.0
	settings.AnnealingFactor = 0.99999

	result, err := hego.Anneal(&initialState, settings)

	if err != nil {
		fmt.Printf("Got error while running Anneal: %v", err)
	} else {
		fmt.Printf("Finished Annealing in %v! Result: %v, Value: %v \n", result.Runtime, result.States[result.Iterations], result.Energies[result.Iterations])
	}
	return
}
