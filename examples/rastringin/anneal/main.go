package main

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/ccssmnn/hego"
)

func rastringin(x, y float64) float64 {
	return 10*2 + (x*x - 10*math.Cos(2*math.Pi*x)) + (y*y - 10*math.Cos(2*math.Pi*y))
}

type state struct {
	v []float64
}

func (s *state) Clone() hego.AnnealState {
	clone := state{v: make([]float64, len(s.v))}
	copy(clone.v, s.v)
	return &clone
}

func (s *state) Neighbor() hego.AnnealState {
	n := state{v: make([]float64, len(s.v))}
	n.v[0] = s.v[0] + rand.NormFloat64()
	n.v[1] = s.v[1] + rand.NormFloat64()
	return &n
}

func (s *state) Energy() float64 {
	return rastringin(s.v[0], s.v[1])
}

func main() {
	initialState := state{v: []float64{1.5, 1.5}}

	settings := hego.AnnealSettings{}
	settings.MaxIterations = 10000
	settings.Verbose = 1000
	settings.Temperature = 10.0
	settings.AnnealingFactor = 0.999

	result, err := hego.Anneal(&initialState, settings)

	if err != nil {
		fmt.Printf("Got error while running Anneal: %v", err)
	} else {
		fmt.Printf("Finished Annealing in %v! Result: %v, Value: %v \n", result.Runtime, result.States[result.Iterations], -result.Energies[result.Iterations])
	}
	return
}
