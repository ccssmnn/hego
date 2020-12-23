package hego

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"text/tabwriter"
	"time"
)

// State represents the current state of the annealing system. Energy is the
// value of the objective function. Neighbor returns another state candidate. Clone
// reproduces this state
type State interface {
	Energy() float64
	Neighbor() State
	Clone() State
}

// AnnealResult represents the result of the Anneal optimization. The last state
// and last energy are the final results. It extends the basic Result type
type AnnealResult struct {
	States   []State
	Energies []float64
	Result
}

// Anneal performs simulated annealing for binary encoded problems such as
// the Knapsack problem
func Anneal(
	initialState State,
	temperature float64,
	annealingFactor float64,
	settings Settings,
) (res AnnealResult, err error) {

	evaluate := func(s State) float64 {
		res.FuncEvaluations++
		return s.Energy()
	}

	start := time.Now()

	state := initialState.Clone()
	energy := evaluate(state)

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

		temperature = temperature * annealingFactor

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
