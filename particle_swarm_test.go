package hego

import (
	"math"
	"math/rand"
	"testing"
)

func TestVerifyPSOSettings(t *testing.T) {
	settings := PSOSettings{}
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
}

func TestPSO(t *testing.T) {
	f := func(x []float64) float64 {
		return math.Pow(x[0], 2)
	}
	init := func() ([]float64, []float64) {
		return []float64{-10 + rand.Float64()*20}, []float64{rand.Float64() * 20.0}
	}
	settings := PSOSettings{}
	_, err := PSO(f, init, settings)
	if err == nil {
		t.Error("PSO should fail with invalid settings")
	}
	settings.MaxIterations = 100
	settings.Verbose = 10
	settings.LearningRate = 1.0
	settings.GlobalWeight = 0.1
	settings.Omega = 0.9
	settings.ParticleWeight = 0.1
	settings.PopulationSize = 10
	res, err := PSO(f, init, settings)
	if err != nil {
		t.Error("PSO should not fail")
	}
	best := res.BestParticles[len(res.BestParticles)-1]
	if best[0] > 0.5 || best[0] < -0.5 {
		t.Errorf("ES algorithm produced unexpected result. Wanted ~0.0, got %v", best[0])
	}
}
