package hego

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"text/tabwriter"
	"time"
)

// AnnealingState represents the current state of the annealing system. Energy is the
// value of the objective function. Neighbor returns another state candidate
type AnnealingState interface {
	Energy() float64
	Neighbor() AnnealingState
}

// SAResult represents the result of the Anneal optimization. The last state
// and last energy are the final results. It extends the basic Result type
type SAResult struct {
	// States hold every state during the process
	States []AnnealingState
	// Energies hold the energy value of every state in the process
	Energies []float64
	Result
}

// SASettings represents the algorithm settings for the simulated annealing
// optimization
type SASettings struct {
	// Temperature is used to determine if another state will be selected or not
	// better states are selected with probability 1, but worse states are selected
	// propability p = exp(state_energy - candidate_energy/temperature)
	// a good value for Temperature is in the range of randomly guessed state energies
	Temperature float64
	// AnnealingFactor is used to decrease the temperature after each iteration
	// When temperature reaches 0, only better states will be accepted which leads
	// to local search / convergence. Thus AnnealingFactor controls after how many
	// iterations convergence might be reached. It's good to reach low temperatures
	// during the last third of iterations
	AnnealingFactor float64
	Settings
}

// Verify returns an error if settings verification fails
func (s *SASettings) Verify() error {
	if s.Temperature <= 0.0 {
		return fmt.Errorf("temperature must be greater that 0.0, got %v", s.Temperature)
	}
	if s.AnnealingFactor > 1.0 || s.AnnealingFactor <= 0.0 {
		return fmt.Errorf("annealing factor must be between 0.0 and 1.0, got %v", s.AnnealingFactor)
	}
	return nil
}

// SA performs simulated annealing algorithm
func SA(
	initialState AnnealingState,
	settings SASettings,
) (res SAResult, err error) {

	err = settings.Verify()
	if err != nil {
		err = fmt.Errorf("settings verification failed: %v", err)
		return
	}

	evaluate := func(s AnnealingState) float64 {
		res.FuncEvaluations++
		return s.Energy()
	}

	start := time.Now()

	state := initialState
	energy := evaluate(state)
	temperature := settings.Temperature
	res.States = append(res.States, state)
	res.Energies = append(res.Energies, energy)
	var buflog bytes.Buffer
	w := tabwriter.NewWriter(
		&buflog, 0, 0, 3, []byte(" ")[0],
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
			state = candidate
			energy = candidateEnergy
		} else {
			probability := math.Exp((energy - candidateEnergy) / temperature)
			if probability > rand.Float64() {
				state = candidate
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
		fmt.Println(buflog.String())
		fmt.Printf("DONE after %v\n", res.Runtime)
	}
	return
}
