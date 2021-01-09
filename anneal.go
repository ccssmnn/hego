package hego

import (
	"fmt"
	"math"
	"math/rand"
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
	// States hold every state during the process (updated on state change)
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

	start := time.Now()

	logger := newLogger("Simulated Annealing", []string{"Iteration", "Temperature", "Energy"}, settings.Verbose, settings.MaxIterations)

	evaluate := func(s AnnealingState) float64 {
		res.FuncEvaluations++
		return s.Energy()
	}

	state := initialState
	energy := evaluate(state)
	temperature := settings.Temperature

	res.States = make([]AnnealingState, 0, settings.MaxIterations)
	res.Energies = make([]float64, settings.MaxIterations)

	for i := 0; i < settings.MaxIterations; i++ {
		candidate := state.Neighbor()
		candidateEnergy := evaluate(candidate)
		update := false
		if candidateEnergy < energy {
			update = true
		} else if math.Exp((energy-candidateEnergy)/temperature) > rand.Float64() {
			update = true
		}
		if update {
			state = candidate
			energy = candidateEnergy
			res.States = append(res.States, state)
		}

		temperature = temperature * settings.AnnealingFactor
		res.Energies[i] = energy
		res.Iterations++
		logger.AddLine(i, []string{
			fmt.Sprint(i),
			fmt.Sprint(temperature),
			fmt.Sprint(energy),
		})
	}

	end := time.Now()
	res.Runtime = end.Sub(start)
	res.Iterations = settings.MaxIterations

	logger.Flush()
	if settings.Verbose > 0 {
		fmt.Printf("Done after %v!\n", res.Runtime)
	}

	return
}
