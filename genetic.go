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

// GeneticSettings represents the settings available in the genetic algorithm
type GeneticSettings struct {
	MutationRate float64
	Elitism      int
	Concurrent   bool
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
		// calculate fitness
		totalFitness := 0.0
		bestFitness := math.MaxFloat64
		worstFitness := -bestFitness
		for idx, can := range pop.candidates {
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
			}
			if fit > worstFitness {
				worstFitness = fit
			}
		}
		res.AveragedFitness = append(res.AveragedFitness, totalFitness/float64(populationSize))
		res.BestFitness = append(res.BestFitness, bestFitness)

		if settings.Elitism > 0 {
			sort.Sort(&pop)
		}

		// select parents
		weights := make([]float64, populationSize)
		for i := 0; i < populationSize; i++ {
			weights[i] = math.Max(worstFitness-pop.candidates[i].fitness, 1e-10)
		}
		parentIds := weightedChoice(weights, populationSize-settings.Elitism)
		parents := make([]GeneticGenome, len(parentIds))
		for index, id := range parentIds {
			parents[index] = pop.candidates[id].genome
		}

		// perform crossover and mutation
		for idx := settings.Elitism; idx < populationSize; idx++ {
			a, b := rand.Intn(len(parents)), rand.Intn(len(parents))
			child := parents[a].Crossover(parents[b])
			if rand.Float64() > settings.MutationRate {
				pop.candidates[idx].genome = child.Mutate()
			} else {
				pop.candidates[idx].genome = child
			}
			pop.candidates[idx].fitness = math.MaxFloat64
		}

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
