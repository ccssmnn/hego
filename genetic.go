package hego

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"
	"text/tabwriter"
	"time"
)

// weightedChoice returns n indizes with a probability defined by weights
// weightedChoice([0.5, 0.3, 0.2], 3) will return 3 indizes. 0 with probability 0.5
// panics if n < 1
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
	indizes := make([]int, 0, n)
	for len(indizes) < n {
		r := rand.Float64() * total
		for i, weight := range weights {
			r -= weight
			if r <= 0.0 {
				indizes = append(indizes, i)
				break
			}
		}
	}
	return indizes
}

// Genome represents a genome (candidate) in the genetic algorithm
// Fitness returns the objective value, Mutate returns a mutated new genome
// and Crossover merges two genomes and returns the child genome
type Genome interface {
	Fitness() float64
	Mutate() Genome
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
	Selection      Selection
	TournamentSize int
	MutationRate   float64
	Elitism        int
	Settings
}

// Verify returns an error, if settings are not valid
func (s *GASettings) Verify() error {
	if s.MaxIterations < 1 {
		return fmt.Errorf("MaxIterations must be greater than 0, not %v", s.MaxIterations)
	}
	if s.MutationRate > 1.0 || s.MutationRate < 0.0 {
		return fmt.Errorf("mutation rate must be between 0.0 and 1.0, not %v", s.MutationRate)
	}
	if s.Elitism < 0 {
		return errors.New("elitism represents the number of best genomes that survive one generation. It cannot be negative")
	}
	if s.Selection == TournamentSelection && s.TournamentSize < 2 {
		return errors.New("When TournamentSelection is set, TournamentSize must be a value above 1")
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
	return weightedChoice(weights, n)
}

func (p population) rankBasedSelection(n int) []int {
	sort.Sort(p)
	weights := make([]float64, len(p))
	for i := range p {
		weights[i] = float64(len(p) - i)
	}
	return weightedChoice(weights, n)
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

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, []byte(" ")[0], tabwriter.AlignRight)

	// logger will log intermediate status data into writer, does nothing when not verbose
	logger := func(i int, avg, best float64) {}
	// flusher will flush writer to stdout if verbose
	flusher := func() {}

	if settings.Verbose > 0 {
		fmt.Println("Starting Genetic Algorithm...")
		fmt.Fprintln(w, "Iteration\tAverage Fitness\tBest Fitness\t")

		logger = func(i int, avg, best float64) {
			if i%settings.Verbose == 0 || i+1 == settings.MaxIterations {
				fmt.Fprintf(w, "%v\t%v\t%v\t\n", i, res.AveragedFitness[i], res.BestFitness[i])
			}
		}

		flusher = func() {
			w.Flush()
			fmt.Printf("DONE after %v\n", res.Runtime)
		}
	}

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
		logger(i, res.AveragedFitness[i], res.BestFitness[i])

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
	flusher()
	res.Runtime = time.Now().Sub(start)
	return res, nil
}
