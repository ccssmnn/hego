package main

import (
	"fmt"
	"math"

	"github.com/ccssmnn/hego"
	"github.com/ccssmnn/hego/mutate"
)

func rastringin(x, y float64) float64 {
	return 10*2 + (x*x - 10*math.Cos(2*math.Pi*x)) + (y*y - 10*math.Cos(2*math.Pi*y))
}

type state struct {
	v []float64
}

func (s *state) Neighbor() hego.AnnealState {
	n := state{v: make([]float64, len(s.v))}
	n.v = mutate.Gauss(s.v, 0.3)
	return &n
}

func (s *state) Energy() float64 {
	return rastringin(s.v[0], s.v[1])
}

func main() {
	initialState := state{v: []float64{5.0, 5.0}}

	settings := hego.AnnealSettings{}
	settings.MaxIterations = 100000
	settings.Verbose = 10000
	settings.Temperature = 10.0
	settings.AnnealingFactor = 0.9999

	result, err := hego.Anneal(&initialState, settings)

	if err != nil {
		fmt.Printf("Got error while running Anneal: %v", err)
	} else {
		fmt.Printf("Finished Annealing in %v! Result: %v, Value: %v \n", result.Runtime, result.States[result.Iterations], result.Energies[result.Iterations])
	}
	return
}
