package hego

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"
	"sync"
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

// GeneticGenome represents a genome (candidate) in the genetic algorithm
// Fitness returns the objective value, Mutate returns a mutated new genome
// and Crossover merges two genomes and returns the child genome
type GeneticGenome interface {
	Fitness() float64
	Mutate() GeneticGenome
	Crossover(other GeneticGenome) GeneticGenome
	GetGene() []interface{}
}

// GeneticResult represents the result of the genetic algorithm
type GeneticResult struct {
	AveragedFitness []float64
	BestFitness     []float64
	BestGenome      []GeneticGenome
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

// GeneticSettings represents the settings available in the genetic algorithm
type GeneticSettings struct {
	Selection      Selection
	TournamentSize int
	MutationRate   float64
	Elitism        int
	Settings
}

// Verify returns an error, if settings are not valid
func (s *GeneticSettings) Verify() error {
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
	genome  GeneticGenome
	fitness float64
}

type population struct {
	candidates []candidate
}

func (p *population) Len() int {
	return len(p.candidates)
}

func (p *population) Less(i, j int) bool {
	return p.candidates[i].fitness < p.candidates[j].fitness
}

func (p *population) Swap(i, j int) {
	p.candidates[i], p.candidates[j] = p.candidates[j], p.candidates[i]
}

func (p *population) fitnessProportionalSelection(n int) []int {
	worst := 0.0
	for _, c := range p.candidates {
		if c.fitness > worst {
			worst = c.fitness
		}
	}
	weights := make([]float64, len(p.candidates))
	for i, c := range p.candidates {
		weights[i] = math.Max(worst-c.fitness, 1e-10)
	}
	return weightedChoice(weights, n)
}

func (p *population) rankBasedSelection(n int) []int {
	sort.Sort(p)
	weights := make([]float64, len(p.candidates))
	for i := range p.candidates {
		weights[i] = float64(len(p.candidates) - i)
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

func (p *population) tournamentSelection(n, size int) []int {
	res := make([]int, n)
	for i := range res {
		// choose tournament candidates from population
		indizes := make([]int, size)
		for j := range indizes {
			indizes[j] = rand.Intn(len(p.candidates))
		}
		// extract fitness from candidates
		weights := make([]float64, size)
		for j, index := range indizes {
			weights[j] = p.candidates[index].fitness
		}
		// determine winner index, which is index in weights slice
		winner := tournament(weights)
		// assign population index to res
		res[i] = indizes[winner]
	}
	return res
}

// Genetic Performs optimization. The optimization follows three steps:
// - for current population calculate fitness
// - select chromosomes with best fitness values with higher propability as parents
// - use parents for reproduction (crossover and mutation)
func Genetic(
	initialPopulation []GeneticGenome,
	settings GeneticSettings,
) (res GeneticResult, err error) {
	populationSize := len(initialPopulation)
	err = settings.Verify()
	if err != nil {
		err = fmt.Errorf("settings verification failed: %v", err)
		return
	}
	// increase FuncEvaluations for every fitness call
	evaluate := func(g GeneticGenome) float64 {
		res.FuncEvaluations++
		return g.Fitness()
	}

	start := time.Now()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, []byte(" ")[0], tabwriter.AlignRight)
	if settings.Verbose > 0 {
		fmt.Println("Starting Genetic Algorithm...")
		fmt.Fprintln(w, "Iteration\tAverage Fitness\tBest Fitness\t")
	}

	pop := population{candidates: make([]candidate, len(initialPopulation))}
	for i := range initialPopulation {
		pop.candidates[i].genome = initialPopulation[i]
	}

	for i := 0; i < settings.MaxIterations; i++ {
		wg := sync.WaitGroup{}
		// calculate fitness
		totalFitness := 0.0
		bestFitness := math.MaxFloat64
		worstFitness := -bestFitness
		bestIndex := -1

		// evalutation of fitness function is independent for each genome
		for idx, can := range pop.candidates {
			wg.Add(1)
			go func(idx int, can candidate) {
				// skip fitness evaluation for elite
				var fit float64
				if i != 0 && idx < settings.Elitism {
					fit = pop.candidates[idx].fitness
				} else {
					fit = evaluate(can.genome)
				}
				totalFitness += fit
				pop.candidates[idx].fitness = fit
				if fit < bestFitness {
					bestFitness = fit
					bestIndex = idx
				}
				if fit > worstFitness {
					worstFitness = fit
				}
				wg.Done()
			}(idx, can)
		}
		wg.Wait()

		res.AveragedFitness = append(res.AveragedFitness, totalFitness/float64(populationSize))
		res.BestFitness = append(res.BestFitness, bestFitness)
		res.BestGenome = append(res.BestGenome, pop.candidates[bestIndex].genome)

		if settings.Elitism > 0 {
			sort.Sort(&pop)
		}

		// select parents
		var parentIds []int
		n := populationSize - settings.Elitism
		switch settings.Selection {
		case RankBasedSelection:
			parentIds = pop.rankBasedSelection(n)
		case TournamentSelection:
			parentIds = pop.tournamentSelection(n, settings.TournamentSize)
		case FitnessProportionalSelection:
			parentIds = pop.fitnessProportionalSelection(n)
		}
		parents := make([]GeneticGenome, len(parentIds))
		for index, id := range parentIds {
			parents[index] = pop.candidates[id].genome
		}

		// perform crossover and mutation
		for idx := settings.Elitism; idx < populationSize; idx++ {
			// crossover and mutation is independent
			wg.Add(1)
			go func(idx int) {
				a, b := rand.Intn(len(parents)), rand.Intn(len(parents))
				child := parents[a].Crossover(parents[b])
				if rand.Float64() > settings.MutationRate {
					pop.candidates[idx].genome = child.Mutate()
				} else {
					pop.candidates[idx].genome = child
				}
				pop.candidates[idx].fitness = math.MaxFloat64
				wg.Done()
			}(idx)
		}
		wg.Wait()

		if settings.Verbose > 0 && (i%settings.Verbose == 0 || i+1 == settings.MaxIterations) {
			fmt.Fprintf(w, "%v\t%v\t%v\t\n", i, res.AveragedFitness[i], res.BestFitness[i])
		}
	}

	if settings.Verbose > 0 {
		w.Flush()
	}

	end := time.Now()
	res.Runtime = end.Sub(start)
	if settings.Verbose > 0 {
		fmt.Printf("DONE after %v\n", res.Runtime)
	}
	res.Iterations = settings.MaxIterations
	return
}
