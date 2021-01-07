package hego

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"text/tabwriter"
	"time"
)

// Ant is the individuum in the population based Ant Colony Optimization (ACO)
type Ant interface {
	// Init initializes the ant for creating a new tour
	Init()
	// Step performs one step to the next position (encoded by int) and returns true when the tour is finished
	Step(next int) bool
	// PerceivePheromone returns a pheromone slice where each element represents the pheromone for the next step (encoded by position)
	PerceivePheromone() []float64
	// DropPheromone leaves pheromone (depending on the performance) on the path of this ant
	DropPheromone(performance float64)
	// Evaporate is called after one iteration and reduces the amount of pheromone on the paths
	Evaporate(factor, min float64)
	// Performance is the objective, lower is better
	Performance() float64
}

// ACOResult represents the result of the ACO
type ACOResult struct {
	// AveragePerformances holds the mean performances for each iteration
	AveragePerformances []float64
	// BestPerformances holds the best performance for each iteration
	BestPerformances []float64
	// BestAnts holds the best Ant for each iteration
	BestAnts []Ant
	Result
}

// ACOSettings represents the settings available in ACO
type ACOSettings struct {
	// Evaporation must be a value in (0, 1] and is used to reduce the amount of pheromone after each iteration
	Evaporation float64
	// MinPheromone is the lowest possible amount of pheromone. Convergence to the true optimum is theoretically only guaranteed for a minpheromone > 0
	MinPheromone float64
	Settings
}

// Verify checks validity of the ACOSettings and returns nil if settings are ok
func (s *ACOSettings) Verify() error {
	if s.Evaporation <= 0 || s.Evaporation > 1 {
		return errors.New("Evaporation must be a value in (0, 1]")
	}
	return nil
}

// ACO performs the ant colony optimization algorithm
func ACO(population []Ant, settings ACOSettings) (res ACOResult, err error) {
	err = settings.Verify()
	if err != nil {
		err = fmt.Errorf("settings verifycation failed: %v", err)
		return
	}

	start := time.Now()

	// increase FuncEvaluations for every performance call
	evaluate := func(a Ant) float64 {
		res.FuncEvaluations++
		return a.Performance()
	}

	var buflog bytes.Buffer
	w := tabwriter.NewWriter(&buflog, 0, 0, 3, []byte(" ")[0], tabwriter.AlignRight)
	// log intermediate status data into writer, does nothing when not verbose
	addLine := func(i int, avg, best float64) {}
	// flush writer to stdout if verbose
	flushTable := func() {}

	if settings.Verbose > 0 {
		fmt.Println("Starting Ant Colony Optimization...")
		fmt.Fprintln(w, "Iteration\tAverage Performance\tBest Performance\t")
		addLine = func(i int, avg, best float64) {
			if i%settings.Verbose == 0 || i+1 == settings.MaxIterations {
				fmt.Fprintf(w, "%v\t%v\t%v\t\n", i, res.AveragePerformances[i], res.BestPerformances[i])
			}
		}
		flushTable = func() {
			w.Flush()
			fmt.Println(buflog.String())
		}
	}

	res.AveragePerformances = make([]float64, settings.MaxIterations)
	res.BestPerformances = make([]float64, settings.MaxIterations)
	res.BestAnts = make([]Ant, settings.MaxIterations)

	for i := 0; i < settings.MaxIterations; i++ {
		totalPerformance := 0.0
		bestPerformance := math.MaxFloat64
		bestIndex := -1
		for antIndex, ant := range population {
			// initialize ant
			ant.Init()
			// create path for this ant
			for {
				options := ant.PerceivePheromone()
				next := weightedChoice(options, 1)[0]
				if ant.Step(next) { // step returns true when ant is done
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
		}
		population[bestIndex].DropPheromone(bestPerformance)
		// population[bestIndex].DropPheromone(bestPerformance)
		population[0].Evaporate(settings.Evaporation, settings.MinPheromone)

		res.AveragePerformances[i] = totalPerformance / float64(len(population))
		res.BestPerformances[i] = bestPerformance
		res.BestAnts[i] = population[bestIndex]
		res.Iterations++
		addLine(i, res.AveragePerformances[i], res.BestPerformances[i])
	}
	end := time.Now()
	res.Runtime = end.Sub(start)
	if settings.Verbose > 0 {
		fmt.Printf("DONE after %v\n", res.Runtime)
	}
	flushTable()
	return
}
