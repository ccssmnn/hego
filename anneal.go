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

// Verify returns an error if settings verification fails
func (s *AnnealSettings) Verify() error {
	if s.MaxIterations <= 0 {
		return fmt.Errorf("iterations must be greater that 0, got %v", s.MaxIterations)
	}
	if s.Temperature <= 0.0 {
		return fmt.Errorf("temperature must be greater that 0.0, got %v", s.Temperature)
	}
	if s.AnnealingFactor > 1.0 || s.AnnealingFactor <= 0.0 {
		return fmt.Errorf("annealing factor must be between 0.0 and 1.0, got %v", s.AnnealingFactor)
	}
	return nil
}

// Anneal performs simulated annealing
func Anneal(
	initialState AnnealState,
	settings AnnealSettings,
) (res AnnealResult, err error) {

	err = settings.Verify()
	if err != nil {
		err = fmt.Errorf("settings verification failed: %v", err)
		return
	}

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

	w := tabwriter.NewWriter(
		os.Stdout, 0, 0, 3, []byte(" ")[0],
		tabwriter.AlignRight,
	)
	if settings.Verbose > 0 {
		fmt.Println("Starting Simulated Annealing...")
		fmt.Fprintln(w, "Iteration\tTemperature\tEnergy\t")
	}

	for i := 0; i < settings.MaxIterations; i++ {
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

		temperature = temperature * settings.AnnealingFactor

		res.States = append(res.States, state)
		res.Energies = append(res.Energies, energy)

		if settings.Verbose > 0 && (i%settings.Verbose == 0 || i+1 == settings.MaxIterations) {
			fmt.Fprintf(w, "%v\t%v\t%v\t\n", i, temperature, energy)
		}
	}

	end := time.Now()
	res.Runtime = end.Sub(start)
	res.Iterations = settings.MaxIterations

	if settings.Verbose > 0 {
		w.Flush()
		fmt.Printf("DONE after %v\n", res.Runtime)
	}

	return
}
