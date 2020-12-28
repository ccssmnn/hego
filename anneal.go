package hego

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"text/tabwriter"
	"time"
)

// AnnealState represents the current state of the annealing system. Energy is the
// value of the objective function. Neighbor returns another state candidate. Clone
// reproduces this state
type AnnealState interface {
	Energy() float64
	Neighbor() AnnealState
	Clone() AnnealState
}

// AnnealResult represents the result of the Anneal optimization. The last state
// and last energy are the final results. It extends the basic Result type
type AnnealResult struct {
	States   []AnnealState
	Energies []float64
	Result
}

// AnnealSettings represents the algorithm settings for the simulated annealing
// optimization
type AnnealSettings struct {
	Temperature     float64
	AnnealingFactor float64
	Settings
}

// Anneal performs simulated annealing for binary encoded problems such as
// the Knapsack problem
func Anneal(
	initialState AnnealState,
	settings AnnealSettings,
) (res AnnealResult, err error) {

	evaluate := func(s AnnealState) float64 {
		res.FuncEvaluations++
		return s.Energy()
	}

	start := time.Now()

	state := initialState.Clone()
	energy := evaluate(state)
	temperature := settings.Temperature

	res.States = append(res.States, state)
	res.Energies = append(res.Energies, energy)
	iter := 0
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, []byte(" ")[0], tabwriter.AlignRight)
	if settings.Verbose > 0 {
		fmt.Println("Starting Simulated Annealing for Bit Problems")
		fmt.Fprintln(w, "Iteration\tTemperature\tEnergy\t")
	}
	for iter < settings.MaxIterations {
		candidate := state.Neighbor()
		candidateEnergy := evaluate(candidate)

		if candidateEnergy < energy {
			state = candidate.Clone()
			energy = candidateEnergy
		} else {
			probability := math.Exp((energy - candidateEnergy) / temperature)
			if probability > rand.Float64() {
				state = candidate.Clone()
				energy = candidateEnergy
			}
		}

		if settings.Verbose > 0 && iter%settings.Verbose == 0 {
			fmt.Fprintf(w, "%v\t%v\t%v\t\n", iter, temperature, energy)
		}

		temperature = temperature * settings.AnnealingFactor

		res.States = append(res.States, state)
		res.Energies = append(res.Energies, energy)

		iter++
	}

	if settings.Verbose > 0 {
		fmt.Fprintf(w, "%v\t%v\t%v\t\n", iter, temperature, energy)
		w.Flush()
	}

	end := time.Now()
	res.Runtime = end.Sub(start)

	res.Iterations = iter
	if settings.Verbose > 0 {
		fmt.Printf("DONE after %v\n", res.Runtime)
	}
	return res, err
}
