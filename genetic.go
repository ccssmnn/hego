package hego

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"
)

// weightedChoice returns n indizes with a probability defined by weights
// weightedChoice([0.5, 0.3, 0.2], 3) will return 3 indizes. 0 with probability 0.5
// panics if n < 1, returns -1 if all weights are 0
func weightedChoice(weights []float64, n int) []int {
	if n < 1 {
		panic("n should be at least 1")
	}
	total := 0.0
	for _, weight := range weights {
		total += weight
	}
	if total == 0.0 {
		return []int{-1}
	}
	indizes := make([]int, n)
	for i := range indizes {
		indizes[i] = -1
		r := rand.Float64() * total
		for j, weight := range weights {
			r -= weight
			if r <= 0.0 {
				indizes[i] = j
				break
			}
		}
	}
	return indizes
}

// binaryWeightedChoice returns n indizes with a probability defined by weights
// it uses a binary search and is more efficient for n > 1 than weightedChoice
// panics if n < 1, returns -1 if all weights are 0
func binaryWeightedChoice(weights []float64, n int) []int {
	if n < 1 {
		panic("number of choices should be 1 or more")
	}
	accumulatedWeights := make([]float64, len(weights))
	cur := 0.0
	for i := 0; i < len(weights); i++ {
		cur += weights[i]
		accumulatedWeights[i] = cur
	}
	if accumulatedWeights[len(accumulatedWeights)-1] == 0.0 {
		return []int{-1}
	}
	makeChoice := func() int {
		target := rand.Float64() * accumulatedWeights[len(weights)-1]
		low, high := 0, len(weights)
		for low < high {
			mid := (low + high) / 2
			distance := accumulatedWeights[mid]
			if distance < target {
				low = mid + 1
			} else if distance > target {
				high = mid
			} else {
				return mid
			}
		}
		return low
	}
	choices := make([]int, n)
	for i := range choices {
		choices[i] = makeChoice()
	}
	return choices
}

// Genome represents a genome (candidate) in the genetic algorithm
type Genome interface {
	// Fitness returns the objective function value for this genome
	Fitness() float64
	// Mutate returns a neighbor of this genome
	Mutate() Genome
	// Crossover merges this and another genome to procude a descendant
	Crossover(other Genome) Genome
}

// GAResult represents the result of the genetic algorithm
type GAResult struct {
	AveragedFitness []float64
	BestFitness     []float64
	BestGenome      []Genome
	Result
}

// Selection encodes different selection variants
type Selection int

const (
	// RankBasedSelection selects parents based on their rank (sorted by fitness)
	RankBasedSelection Selection = iota
	// TournamentSelection performs a tournament of randomly selected genomes
	// and selects the winner
	TournamentSelection
	// FitnessProportionalSelection determines the chance of a genome to be
	// selected by its fitness value compared to the total fitness of the population
	FitnessProportionalSelection
)

// GASettings represents the settings available in the genetic algorithm
type GASettings struct {
	// Selection defines the type of selection to be used
	Selection Selection
	// TournamentSize defines the size of a tournament (only necessary for TournamentSelection)
	TournamentSize int
	// MutationRate is the probability of a candidate to mutate after crossover
	MutationRate float64
	// Elitism is the number of best candidates to pass over to the next generation without selection
	Elitism int
	Settings
}

// Verify returns an error, if settings are not valid
func (s *GASettings) Verify() error {
	if s.MutationRate > 1.0 || s.MutationRate < 0.0 {
		return fmt.Errorf("mutation rate must be between 0.0 and 1.0, not %v", s.MutationRate)
	}
	if s.Elitism < 0 {
		return errors.New("elitism cannot be negative")
	}
	if s.Selection == TournamentSelection && s.TournamentSize < 2 {
		return errors.New("when TournamentSelection is set, TournamentSize must be a value above 1")
	}
	return nil
}

type candidate struct {
	genome  Genome
	fitness float64
}

type population []candidate

func (p population) Len() int {
	return len(p)
}

func (p population) Less(i, j int) bool {
	return p[i].fitness < p[j].fitness
}

func (p population) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p population) fitnessProportionalSelection(n int) []int {
	worst := 0.0
	for _, c := range p {
		if c.fitness > worst {
			worst = c.fitness
		}
	}
	weights := make([]float64, len(p))
	for i, c := range p {
		weights[i] = math.Max(worst-c.fitness, 1e-10)
	}
	return binaryWeightedChoice(weights, n)
}

func (p population) rankBasedSelection(n int) []int {
	sort.Sort(p)
	weights := make([]float64, len(p))
	for i := range p {
		weights[i] = float64(len(p) - i)
	}
	return binaryWeightedChoice(weights, n)
}

func tournament(weights []float64) int {
	contesters := make([]int, len(weights))
	for i := range contesters {
		contesters[i] = i
	}
	for len(contesters) > 1 {
		winners := make([]int, len(contesters)/2)
		for i := range winners {
			a, b := contesters[2*i], contesters[2*i+1]
			// lower is better! we are minimizing
			if weights[a] > weights[b] {
				winners[i] = b
			} else {
				winners[i] = a
			}
		}
		contesters = winners
	}
	return contesters[0]
}

func (p population) tournamentSelection(n, size int) []int {
	res := make([]int, n)
	for i := range res {
		// choose tournament candidates from population
		indizes := make([]int, size)
		for j := range indizes {
			indizes[j] = rand.Intn(len(p))
		}
		// extract fitness from candidates
		weights := make([]float64, size)
		for j, index := range indizes {
			weights[j] = p[index].fitness
		}
		// determine winner index, which is index in weights slice
		winner := tournament(weights)
		// assign population index to res
		res[i] = indizes[winner]
	}
	return res
}

func (p population) selectParents(settings *GASettings) []int {
	n := len(p) - settings.Elitism
	var parentIds []int
	switch settings.Selection {
	case RankBasedSelection:
		parentIds = p.rankBasedSelection(n)
	case TournamentSelection:
		parentIds = p.tournamentSelection(n, settings.TournamentSize)
	case FitnessProportionalSelection:
		parentIds = p.fitnessProportionalSelection(n)
	}
	return parentIds
}

// GA Performs optimization. The optimization follows three steps:
// - for current population calculate fitness
// - select chromosomes with best fitness values with higher propability as parents
// - use parents for reproduction (crossover and mutation)
func GA(
	initialPopulation []Genome,
	settings GASettings,
) (res GAResult, err error) {
	err = settings.Verify()
	if err != nil {
		err = fmt.Errorf("settings verification failed: %v", err)
		return
	}
	// increase FuncEvaluations for every fitness call
	evaluate := func(g Genome) float64 {
		res.FuncEvaluations++
		return g.Fitness()
	}

	start := time.Now()
	logger := newLogger("Genetic Algorithm", []string{"Iteration", "Average Fitness", "Best Fitness"}, settings.Verbose, settings.MaxIterations)

	pop := make(population, len(initialPopulation))
	for i := range initialPopulation {
		pop[i].genome = initialPopulation[i]
		pop[i].fitness = evaluate(pop[i].genome)
	}

	res.AveragedFitness = make([]float64, settings.MaxIterations)
	res.BestFitness = make([]float64, settings.MaxIterations)
	res.BestGenome = make([]Genome, settings.MaxIterations)

	for i := 0; i < settings.MaxIterations; i++ {
		// FITNESS EVALUATION
		totalFitness := 0.0
		bestFitness := math.MaxFloat64
		bestIndex := -1
		for idx, g := range pop {
			totalFitness += g.fitness
			if g.fitness < bestFitness {
				bestFitness = g.fitness
				bestIndex = idx
			}
		}
		res.AveragedFitness[i] = totalFitness / float64(len(pop))
		res.BestFitness[i] = bestFitness
		res.BestGenome[i] = pop[bestIndex].genome

		logger.AddLine(i, []string{
			fmt.Sprint(i),
			fmt.Sprint(res.AveragedFitness[i]),
			fmt.Sprint(res.BestFitness[i]),
		})

		// SELECTION
		parentIds := pop.selectParents(&settings)

		// CROSSOVER & MUTATION
		// TODO: for elitism << len(pop) it is more efficient to extract smallest n instead of sorting
		if settings.Elitism > 0 {
			sort.Sort(&pop)
		}
		for idx := settings.Elitism; idx < len(pop); idx++ {
			parent1 := pop[parentIds[rand.Intn(len(parentIds))]].genome
			parent2 := pop[parentIds[rand.Intn(len(parentIds))]].genome
			if rand.Float64() < settings.MutationRate {
				pop[idx].genome = parent1.Crossover(parent2).Mutate()
			} else {
				pop[idx].genome = parent1.Crossover(parent2)
			}
			pop[idx].fitness = evaluate(pop[idx].genome)
		}
		res.Iterations++
	}
	logger.Flush()
	res.Runtime = time.Now().Sub(start)
	if settings.Verbose > 0 {
		fmt.Printf("DONE after %v\n", res.Runtime)
	}
	return res, nil
}
