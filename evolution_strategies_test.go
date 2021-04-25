package hego

import (
	"testing"
)

func TestVerifyESSettings(t *testing.T) {
	settings := ESSettings{}
	err := settings.Verify()
	if err == nil {
		t.Error("expected es settings verification to fail with no custom settings")
	}
	settings.LearningRate = 0.1
	err = settings.Verify()
	if err == nil {
		t.Error("expected es settings verification to fail without population size set")
	}
	settings.PopulationSize = 10
	err = settings.Verify()
	if err == nil {
		t.Error("expected es settings verification to fail without sigma set")
	}
	settings.NoiseSigma = 0.5
	err = settings.Verify()
	if err != nil {
		t.Error("settings verification should pass with valid lr, size, sigma")
	}
}

func TestES(t *testing.T) {
	f := func(x []float64) float64 {
		return x[0] * x[0]
	}
	x0 := []float64{10.0}
	settings := ESSettings{}
	_, err := ES(f, x0, settings)
	if err == nil {
		t.Error("ES should fail with invalid settings")
	}
	settings.MaxIterations = 10
	settings.Verbose = 10
	settings.LearningRate = 1.0
	settings.NoiseSigma = 0.1
	settings.PopulationSize = 10
	res, err := ES(f, x0, settings)
	if err != nil {
		t.Errorf("Unexpected error in ES algorithm: %v", err)
	}
	best := res.BestCandidate
	if best[0] > 0.5 || best[0] < -0.5 {
		t.Errorf("ES algorithm produced unexpected result. Wanted ~0.0, got %v", best[0])
	}
}
