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

func (s *state) Clone() hego.State {
	clone := state{v: make([]float64, len(s.v))}
	copy(clone.v, s.v)
	return &clone
}

func (s *state) Neighbor() hego.State {
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
	settings := hego.Settings{
		MaxIterations: 10000,
		Verbose:       1000,
	}

	temperature := 10.0
	annealingFactor := 0.999

	result, err := hego.Anneal(&initialState, temperature, annealingFactor, settings)

	if err != nil {
		fmt.Printf("Got error while running Anneal: %v", err)
	} else {
		fmt.Printf("Finished Annealing in %v! Result: %v, Value: %v \n", result.Runtime, result.States[result.Iterations], -result.Energies[result.Iterations])
	}
	return
}
