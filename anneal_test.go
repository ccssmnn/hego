package hego

import (
	"math"
	"math/rand"
	"testing"
)

type state float64

func (b state) Energy() float64 {
	return float64(b * b)
}

func (b state) Neighbor() AnnealingState {
	return b + state(rand.NormFloat64()*0.1)
}

func TestVerify(t *testing.T) {
	settings := SASettings{}
	settings.Temperature = 100.0
	settings.AnnealingFactor = 0.9
	err := settings.Verify()
	if err != nil {
		t.Error("verification should pass for temperature above 0 and annealing factor in (0,1]")
	}
	settings.Temperature = 0.0
	err = settings.Verify()
	if err == nil {
		t.Error("verification should fail for temperature = 0.0")
	}
	settings.Temperature = 100.0
	settings.AnnealingFactor = 1.1
	err = settings.Verify()
	if err == nil {
		t.Error("verification should fail for annealing factor above 1")
	}
	settings.AnnealingFactor = 0.0
	err = settings.Verify()
	if err == nil {
		t.Error("verification should fail for annealing factor <= 0")
	}
}

// TestAnnealBit runs the AnnealBit method and checks for errors
func TestSA(t *testing.T) {
	initialState := state(20.0)

	settings := SASettings{}
	res, err := SA(initialState, settings)
	if err == nil {
		t.Error("SA should fail with invalid settings")
	}
	settings.Temperature = 50.0
	settings.AnnealingFactor = 0.99
	settings.MaxIterations = 1000
	settings.Verbose = 1
	settings.KeepHistory = true
	res, err = SA(initialState, settings)
	if err != nil {
		t.Errorf("Error while running Anneal main algorithm: %v", err)
	}
	if len(res.Energies) == 0 {
		t.Error("Energies list should not be empty")
	}
	if math.Abs(res.Energy) > 0.5 {
		t.Error("unexpected solution")
	}
}
