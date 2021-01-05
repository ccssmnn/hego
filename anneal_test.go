package hego

import (
	"math/rand"
	"testing"
)

type state []bool

func (b state) Energy() float64 {
	return rand.Float64()
}

func (b state) Neighbor() AnnealingState {
	return b
}

// TestAnnealBit runs the AnnealBit method and checks for errors
func TestSA(t *testing.T) {
	initialState := state{false, true, false}

	settings := SASettings{}
	settings.Temperature = 100.0
	settings.AnnealingFactor = 0.9
	settings.MaxIterations = 10
	settings.Verbose = 0

	res, err := SA(initialState, settings)
	if err != nil {
		t.Errorf("Error while running Anneal main algorithm: %v", err)
	}
	// maxiterations + 1 because the initial state is not included in the counting
	if len(res.States) != settings.MaxIterations+1 {
		t.Errorf("Wrong number of states received: expected %v, got %v", settings.MaxIterations, len(res.States))
	}
}
