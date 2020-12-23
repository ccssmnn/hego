package hego

import (
	"math/rand"
	"testing"
)

type bitState struct {
	state []bool
}

func (b *bitState) Clone() State {
	clone := bitState{state: make([]bool, len(b.state))}
	copy(clone.state, b.state)
	return &clone
}

func (b *bitState) Energy() float64 {
	return rand.Float64()
}

func (b *bitState) Neighbor() State {
	n := bitState{state: make([]bool, len(b.state))}
	for i := range n.state {
		n.state[i] = rand.Float64() < 0.5
	}
	return &n
}

// TestAnnealBit runs the AnnealBit method and checks for errors
func TestAnnealBit(t *testing.T) {
	initialState := bitState{state: []bool{false, true, false}}
	temperature := 100.0
	annealingFactor := 0.9
	settings := Settings{
		MaxIterations: 10,
		Verbose:       1,
	}
	res, err := Anneal(&initialState, temperature, annealingFactor, settings)
	if err != nil {
		t.Errorf("Error while running Anneal main algorithm: %v", err)
	}
	// maxiterations + 1 because the initial state is not included in the counting
	if len(res.States) != settings.MaxIterations+1 {
		t.Errorf("Wrong number of states received: expected %v, got %v", settings.MaxIterations, len(res.States))
	}
}
