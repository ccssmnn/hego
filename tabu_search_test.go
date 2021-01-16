package hego

import (
	"math/rand"
	"testing"
)

type tabuState []bool

func (b tabuState) Equal(other TabuState) bool {
	return false
}

func (b tabuState) Objective() float64 {
	return rand.Float64()
}

func (b tabuState) Neighbor() TabuState {
	return b
}

func TestVerifyTSSettings(t *testing.T) {
	settings := TSSettings{}
	settings.NeighborhoodSize = 5
	settings.TabuListSize = 5
	err := settings.Verify()
	if err != nil {
		t.Error("verification should pass for NeighborhoodSize above 0 and TabuListSize above 1")
	}
	settings.NeighborhoodSize = 0
	err = settings.Verify()
	if err == nil {
		t.Error("verification should fail for NeighborhoodSize = 0")
	}
	settings.NeighborhoodSize = 5
	settings.TabuListSize = 0
	err = settings.Verify()
	if err == nil {
		t.Error("verification should fail for tabu list size below 1")
	}
}

// TestAnnealBit runs the AnnealBit method and checks for errors
func TestTS(t *testing.T) {
	initialState := tabuState{false, true, false}

	settings := TSSettings{}
	_, err := TS(initialState, settings)
	if err == nil {
		t.Error("SA should fail with invalid settings")
	}
	settings.NeighborhoodSize = 100.0
	settings.TabuListSize = 5
	settings.MaxIterations = 10
	settings.Verbose = 1
	_, err = TS(initialState, settings)
	if err != nil {
		t.Errorf("Error while running Anneal main algorithm: %v", err)
	}
}
