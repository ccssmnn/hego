package hego

import (
	"math/rand"
	"testing"
)

type genome struct {
	state []bool
}

func (b *genome) Crossover(other GeneticGenome) GeneticGenome {
	clone := genome{state: make([]bool, len(b.state))}
	copy(clone.state, b.state)
	return &clone
}

func (b *genome) Fitness() float64 {
	return rand.Float64()
}

func (b *genome) Mutate() GeneticGenome {
	n := genome{state: make([]bool, len(b.state))}
	for i := range n.state {
		n.state[i] = rand.Float64() < 0.5
	}
	return &n
}

func (b *genome) GetGene() []interface{} {
	gene := make([]interface{}, len(b.state))
	for i, value := range b.state {
		gene[i] = value
	}
	return gene
}

func TestGenetic(t *testing.T) {
	population := make([]GeneticGenome, 0)
	for i := 0; i < 10; i++ {
		candidate := genome{}
		for index := range candidate.state {
			candidate.state[index] = rand.Float64() > 0.5
		}
		population = append(population, &candidate)
	}

	settings := GeneticSettings{}
	settings.MutationRate = 0.5
	settings.Elitism = 1
	settings.MaxIterations = 10
	settings.Verbose = 0

	res, err := Genetic(population, settings)
	if err != nil {
		t.Errorf("Error while running Anneal main algorithm: %v", err)
	}
	if len(res.BestFitness) != settings.MaxIterations {
		t.Errorf("Wrong number of states received: expected %v, got %v", settings.MaxIterations, len(res.BestFitness))
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
