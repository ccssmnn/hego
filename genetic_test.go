package hego

import (
	"math/rand"
	"testing"
)

var crossoverCount int
var fitnessCount int
var mutateCount int

type genome struct {
	state [10]bool
}

func (b genome) Crossover(other Genome) Genome {
	crossoverCount++
	return b
}

func (b genome) Fitness() float64 {
	fitnessCount++
	return 1.0
}

func (b genome) Mutate() Genome {
	mutateCount++
	return b
}

func TestGenetic(t *testing.T) {

	settings := GASettings{}
	settings.MutationRate = 0.0
	settings.Elitism = 0
	settings.MaxIterations = 10
	settings.Verbose = 0
	populationSize := 10

	population := make([]Genome, populationSize)
	for i := range population {
		candidate := genome{}
		for index := range candidate.state {
			candidate.state[index] = rand.Float64() > 0.5
		}
		population[i] = candidate
	}

	crossoverCount = 0
	fitnessCount = 0
	mutateCount = 0

	res, err := GA(population, settings)

	if err != nil {
		t.Errorf("Error while running Anneal main algorithm: %v", err)
	}
	if res.Iterations != settings.MaxIterations {
		t.Errorf("unexpected number of iterations. Expected %v, got %v", settings.MaxIterations, res.Iterations)
	}
	expectedCrossoverCount := settings.MaxIterations*populationSize - settings.Elitism
	if crossoverCount != expectedCrossoverCount {
		t.Errorf("unexpected number of crossover operations: Expected %v, got %v", expectedCrossoverCount, crossoverCount)
	}
	expectedFitnessCount := settings.MaxIterations*populationSize - settings.Elitism
	if crossoverCount != expectedFitnessCount {
		t.Errorf("unexpected number of fitness calls: Expected %v, got %v", expectedFitnessCount, fitnessCount)
	}
	if mutateCount != 0 {
		t.Errorf("unexpected number of mutate operations: Expected %v, got %v", 0, mutateCount)
	}

	// retry with 100% mutation rate
	settings.MutationRate = 1.0
	res, err = GA(population, settings)
	expectedMutateCount := settings.MaxIterations*populationSize - settings.Elitism
	if mutateCount != expectedCrossoverCount {
		t.Errorf("unexpected number of mutate operations: Expected %v, got %v", expectedMutateCount, mutateCount)
	}
}

func TestWeightedChoice(t *testing.T) {
	weights := []float64{1.0, 2.0, 0.0}
	n := 20
	choices := weightedChoice(weights, n)
	if len(choices) != n {
		t.Errorf("expected number of choices to be %v, got %v", n, len(choices))
	}
	for _, choice := range choices {
		if choice == 2 {
			t.Error("2 should not be a choice")
		}
	}
}

func TestTournament(t *testing.T) {
	weights := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0}
	index := tournament(weights)
	if index != 0 {
		t.Errorf("expected index 0 to win the tournament, got %v", index)
	}
}
