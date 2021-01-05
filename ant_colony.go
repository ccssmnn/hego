package hego

import (
	"errors"
	"fmt"
	"math"
	"os"
	"text/tabwriter"
	"time"
)

// Ant is the individuum in the population based Ant Colony Optimization (ACO)
type Ant interface {
	Init()
	Step(next int) bool
	PerceivePheromone() []float64
	DropPheromone(performance float64)
	Evaporate(factor, min float64)
	Performance() float64
}

// ACOResult represents the result of the ACO
type ACOResult struct {
	AveragePerformances []float64
	BestPerformances    []float64
	BestAnts            []Ant
	Result
}

// ACOSettings represents the settings available in ACO
type ACOSettings struct {
	Evaporation  float64
	MinPheromone float64
	Settings
}

// Verify checks validity of the ACOSettings and returns nil if everything is fine
func (s *ACOSettings) Verify() error {
	if s.Evaporation <= 0 || s.Evaporation > 1 {
		return errors.New("Evaporation must be a value in (0, 1]")
	}
	return nil
}

// InitializePheromoneMatrix generates a pheromone matrix based on a distance
// matrix and a starter value
func InitializePheromoneMatrix(dim int, starter float64) [][]float64 {
	pheromones := make([][]float64, 0, dim)
	for i := 0; i < dim; i++ {
		line := make([]float64, dim)
		for j := 0; j < dim; j++ {
			line[j] = starter
		}
		pheromones = append(pheromones, line)
	}
	return pheromones
}

// ACO performs the ant colony optimization algorithm
func ACO(population []Ant, settings ACOSettings) (res ACOResult, err error) {
	err = settings.Verify()
	if err != nil {
		err = fmt.Errorf("settings verifycation failed: %v", err)
		return
	}

	// increase FuncEvaluations for every performance call
	evaluate := func(a Ant) float64 {
		res.FuncEvaluations++
		return a.Performance()
	}

	start := time.Now()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, []byte(" ")[0], tabwriter.AlignRight)
	if settings.Verbose > 0 {
		fmt.Println("Starting ACO Algorithm...")
		fmt.Fprintln(w, "Iteration\tAverage Performance\tBest Performance\t")
	}
	for i := 0; i < settings.MaxIterations; i++ {
		totalPerformance := 0.0
		bestPerformance := math.MaxFloat64
		worstPerformance := -bestPerformance
		bestIndex := -1
		for antIndex, ant := range population {
			// initialize ant
			ant.Init()
			// create path for this ant
			for {
				options := ant.PerceivePheromone()
				next := weightedChoice(options, 1)[0]
				if next == -1 {
					break
				}
				done := ant.Step(next)
				if done {
					break
				}
			}
			// evaluate path
			performance := evaluate(ant)
			totalPerformance += performance
			if performance < bestPerformance {
				bestPerformance = performance
				bestIndex = antIndex
			}
			if performance > worstPerformance {
				worstPerformance = performance
			}
		}
		population[bestIndex].DropPheromone(bestPerformance)
		population[0].Evaporate(settings.Evaporation, settings.MinPheromone)

		res.AveragePerformances = append(res.AveragePerformances, totalPerformance/float64(len(population)))
		res.BestPerformances = append(res.BestPerformances, bestPerformance)
		res.BestAnts = append(res.BestAnts, population[bestIndex])

		if settings.Verbose > 0 && (i%settings.Verbose == 0 || i+1 == settings.MaxIterations) {
			fmt.Fprintf(w, "%v\t%v\t%v\t\n", i, res.AveragePerformances[i], res.BestPerformances[i])
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
