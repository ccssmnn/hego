package hego

import (
	"math"
	"math/rand"
	"testing"
)

type tabuState float64

func (b tabuState) Equal(other TabuState) bool {
	return b == other.(tabuState)
}

func (b tabuState) Objective() float64 {
	return float64(b * b)
}

func (b tabuState) Neighbor() TabuState {
	return b + tabuState(rand.NormFloat64()*0.1)
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
	initialState := tabuState(10.0)

	settings := TSSettings{}
	_, err := TS(initialState, settings)
	if err == nil {
		t.Error("SA should fail with invalid settings")
	}
	settings.NeighborhoodSize = 100
	settings.TabuListSize = 5
	settings.MaxIterations = 100
	settings.Verbose = 1
	res, err := TS(initialState, settings)
	if err != nil {
		t.Errorf("Error while running tabu search algorithm: %v", err)
	}
	if math.Abs(float64(res.BestState.(tabuState))) > 0.5 {
		t.Errorf("Unexpected optimization Result")
	}
	if len(res.States) != 0 {
		t.Error("states should be empty")
	}
	settings.Verbose = 0
	settings.KeepHistory = true
	res, err = TS(initialState, settings)
	if err != nil {
		t.Errorf("Error while running tabu search algorithm: %v", err)
	}
	if len(res.States) == 0 {
		t.Error("states should not be empty")
	}
}
