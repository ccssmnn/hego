package hego

import (
	"fmt"
	"math"
	"math/rand"
	"os"
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
	return nil
}

// Genetic Performs optimization. The optimization follows three steps:
// - for current population calculate fitness
// - select chromosomes with best fitness values with higher propability as parents
// - use parents for reproduction (crossover and mutation)
func Genetic(
	population []GeneticGenome,
	settings GeneticSettings,
) (res GeneticResult, err error) {

	err = settings.Verify()
	if err != nil {
		err = fmt.Errorf("settings verification failed: %v", err)
		return
	}

	start := time.Now()

	evaluate := func(g GeneticGenome) float64 {
		res.FuncEvaluations++
		return g.Fitness()
	}

	fitnesses := make([]float64, len(population))

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, []byte(" ")[0], tabwriter.AlignRight)
	if settings.Verbose > 0 {
		fmt.Println("Starting Genetic Algorithm...")
		fmt.Fprintln(w, "Iteration\tAverage Fitness\tBest Fitness\t")
	}

	for i := 0; i < settings.MaxIterations; i++ {
		// calculate fitness
		totalFitness := 0.0
		bestFitness := math.MaxFloat64
		worstFitness := -bestFitness
		for idx, genome := range population {
			fit := evaluate(genome)
			totalFitness += fit
			fitnesses[idx] = fit
			if fit < bestFitness {
				bestFitness = fit
			}
			if fit > worstFitness {
				worstFitness = fit
			}
		}
		res.AveragedFitness = append(res.AveragedFitness, totalFitness/float64(len(population)))
		res.BestFitness = append(res.BestFitness, bestFitness)

		// select parents
		weights := make([]float64, len(fitnesses))
		for i, fit := range fitnesses {
			weights[i] = math.Max(worstFitness-fit, 1e-10)
		}
		parentIds := weightedChoice(weights, len(fitnesses))
		parents := make([]GeneticGenome, len(parentIds))
		for index, id := range parentIds {
			parents[index] = population[id]
		}

		// perform crossover and mutation
		for idx := range population {
			a, b := rand.Intn(len(parents)), rand.Intn(len(parents))
			population[idx] = parents[a].Crossover(parents[b])
			if rand.Float64() > settings.MutationRate {
				population[idx] = population[idx].Mutate()
			}
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
