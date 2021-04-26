package hego

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"
)

// ESResult represents the result of the evolution strategy algorithm
type ESResult struct {
	Candidates        [][]float64
	AverageObjectives []float64
	BestObjectives    []float64
	BestCandidate     []float64
	BestObjective     float64
	Result
}

// ESSettings represents settings for the evolution strategy algorithm
type ESSettings struct {
	// PopulationSize is the number of noise vectors to create for each iteration
	// these noise vectors are used to create a gradient estimate, so population size should not
	// be too small
	PopulationSize int
	// LearningRate is the factor to determine the step size after each iteration
	// a step is made by calculating x = x + learningRate * gradient_estimate(x)
	LearningRate float64
	// NoiseSigma is the sigma value for noise generated. A higher sigma results in a wider
	// search spread, but might result in inaccuracies for the gradient estimate
	NoiseSigma float64
	Settings
}

// Verify checks the validity of the settings and returns nil if everything is ok
func (s *ESSettings) Verify() error {
	if s.LearningRate <= 0.0 {
		return fmt.Errorf("learning rate must be a value above 0.0, got %v", s.LearningRate)
	}
	if s.PopulationSize <= 1 {
		return fmt.Errorf("population size must be greater than 1, got %v", s.PopulationSize)
	}
	if s.NoiseSigma == 0.0 {
		return errors.New("sigma = 0.0 leads to no search at all")
	}
	return nil
}

// ES performs Evolutionary Strategy algorithm suited for minimizing
// a real valued function (objective) from a starting point x0
// It takes advantage of population based gradient updates, where each iteration a population
// is generated from noise added to the current and used to estimate the gradient.
func ES(
	objective func(x []float64) float64,
	x0 []float64,
	settings ESSettings) (res ESResult, err error) {
	err = settings.Verify()
	if err != nil {
		err = fmt.Errorf("settings verification failed: %v", err)
		return res, err
	}
	start := time.Now()
	logger := newLogger("Evolution Strategy Algorithm", []string{"Iteration", "Population Mean", "Current Candidate"}, settings.Verbose, settings.MaxIterations)
	// increase funcEvaluations counter for every call to objective
	evaluate := func(x []float64) float64 {
		res.FuncEvaluations++
		return objective(x)
	}
	// write noise into x
	initNoise := func(x []float64) {
		for i := range x {
			x[i] = rand.NormFloat64() * settings.NoiseSigma
		}
	}
	// add x to noise vector
	combineWithNoise := func(noise, x []float64) {
		for i := range noise {
			noise[i] += x[i]
		}
	}

	candidate := make([]float64, len(x0))
	copy(candidate, x0)

	if settings.KeepHistory {
		res.BestObjectives = make([]float64, settings.MaxIterations)
		res.AverageObjectives = make([]float64, settings.MaxIterations)
		res.Candidates = make([][]float64, settings.MaxIterations)
	}

	res.BestObjective = math.MaxFloat64
	res.BestCandidate = make([]float64, len(x0))

	// initialize memory for population and their rewards
	population := make([][]float64, settings.PopulationSize)
	for i := range population {
		population[i] = make([]float64, len(x0))
	}
	rewards := make([]float64, settings.PopulationSize)

	for i := 0; i < settings.MaxIterations; i++ {

		totalReward := 0.0
		bestReward := math.MaxFloat64
		for j := range population {
			// create new candidate with noise
			initNoise(population[j])
			combineWithNoise(population[j], candidate)
			reward := evaluate(population[j])
			rewards[j] = reward
			totalReward += reward
			if reward < bestReward {
				bestReward = reward
			}
		}

		// compute standart deviation of rewards
		meanReward := totalReward / float64(settings.PopulationSize)
		stdDev := 0.0
		for _, reward := range rewards {
			stdDev += math.Pow(reward-meanReward, 2)
		}
		stdDev = math.Sqrt(stdDev / float64(settings.PopulationSize))

		// update candidate
		for j := range candidate {
			// estimate gradient
			gradientEstimate := 0.0
			for index, individuum := range population {
				gradientEstimate += individuum[j] * (rewards[index] - meanReward) / stdDev
			}
			gradientEstimate *= 1.0 / (float64(settings.PopulationSize) * settings.NoiseSigma)

			// perform gradient step towards minimum
			candidate[j] -= settings.LearningRate * gradientEstimate
		}
		// update result
		if settings.KeepHistory {
			res.Candidates[i] = make([]float64, len(candidate))
			copy(res.Candidates[i], candidate)
			res.BestObjectives[i] = bestReward
			res.AverageObjectives[i] = meanReward
		}
		if res.BestObjective > bestReward {
			res.BestObjective = bestReward
			copy(res.BestCandidate, candidate)
		}

		logger.AddLine(i, []string{
			fmt.Sprint(i),
			fmt.Sprint(meanReward),
			fmt.Sprint(objective(candidate)),
		})
	}

	end := time.Now()
	res.Runtime = end.Sub(start)
	res.Iterations = settings.MaxIterations
	logger.Flush()
	if settings.Verbose > 0 {
		fmt.Printf("Done after %v!\n", res.Runtime)
	}
	return res, nil
}
