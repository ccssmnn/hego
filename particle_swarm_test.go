package hego

import (
	"math/rand"
	"testing"
)

func TestVerifyPSOSettings(t *testing.T) {
	settings := PSOSettings{}
	settings.PopulationSize = 10
	settings.LearningRate = 0.1
	settings.Omega = 0.1
	settings.GlobalWeight = 0.1
	settings.ParticleWeight = 0.1

	err := settings.Verify()
	if err != nil {
		t.Error("expected verification to pass with these valid settings")
	}

	settings.PopulationSize = 0
	err = settings.Verify()
	if err == nil {
		t.Error("expected verification to fail with populationsize 0")
	}

	settings.PopulationSize = 10
	settings.LearningRate = 0.0
	err = settings.Verify()
	if err == nil {
		t.Error("expected verification to fail with learningrate 0")
	}

	settings.LearningRate = 0.0
	settings.Omega = -1.0
	err = settings.Verify()
	if err == nil {
		t.Error("expected verification to fail with negative omega")
	}

	settings.Omega = 1.0
	settings.GlobalWeight = -1.0
	err = settings.Verify()
	if err == nil {
		t.Error("expected verification to fail with negative globalweight")
	}
	settings.GlobalWeight = 1.0
	settings.ParticleWeight = -1.0
	err = settings.Verify()
	if err == nil {
		t.Error("expected verification to fail with negative ParticleWeight")
	}
	settings.ParticleWeight = 0.0
	settings.GlobalWeight = 0.0
	err = settings.Verify()
	if err == nil {
		t.Error("expected verification to fail with both zero particle weight and global weight")
	}
}

func TestPSO(t *testing.T) {
	f := func(x []float64) float64 {
		return x[0] * x[0]
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
