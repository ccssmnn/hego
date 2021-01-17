package hego

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"
)

// PSOResult represents the results of the particle swarm optimization
type PSOResult struct {
	BestParticles  [][]float64
	BestObjectives []float64
	Result
}

// PSOSettings represents settings for the particle swarm optimization
type PSOSettings struct {
	// PopulationSize is the number of particles
	PopulationSize int
	// LearningRate determines the movement size of each particle
	LearningRate float64
	// Omega is the weight of the current velocity, a momentum
	Omega float64
	// GlobalWeight determines how much a particle should drift towards the global optimum
	GlobalWeight float64
	// ParticleWeight determines how much a particle should drift towards the best known position of this particle
	ParticleWeight float64
	Settings
}

// Verify checks the validity of the settings and returns nil if everything is ok
func (s *PSOSettings) Verify() error {
	if s.PopulationSize <= 1 {
		return fmt.Errorf("population size must be greater than 1, got %v", s.PopulationSize)
	}
	if s.LearningRate <= 0.0 {
		return fmt.Errorf("learning rate must be greater than 0, got %v", s.LearningRate)
	}
	if s.Omega < 0.0 {
		return fmt.Errorf("omega should not be smaller than 0, got %v", s.Omega)
	}
	if s.GlobalWeight < 0.0 {
		return fmt.Errorf("GlobalWeight should not be smaller than 0, got %v", s.GlobalWeight)
	}
	if s.ParticleWeight < 0.0 {
		return fmt.Errorf("ParticleWeight should not be smaller than 0, got %v", s.ParticleWeight)
	}
	if s.ParticleWeight == 0.0 && s.GlobalWeight == 0.0 {
		return errors.New("when ParticleWeight and GlobalWeight are set to 0, the velocity will not change at all")
	}
	return nil
}

// PSO performs particle swarm optimization. Objective is the function to minimize, init initializes a tupe of particle and velocity, settings holds algorithm settings
func PSO(
	objective func(x []float64) float64,
	init func() ([]float64, []float64),
	settings PSOSettings) (res PSOResult, err error) {
	err = settings.Verify()
	if err != nil {
		err = fmt.Errorf("settings verification failed: %v", err)
		return res, err
	}
	start := time.Now()
	logger := newLogger("Particle Swarm Optimization", []string{"Iteration", "Population Mean", "Population Best"}, settings.Verbose, settings.MaxIterations)
	// increase funcEvaluations counter for every call to objective
	evaluate := func(x []float64) float64 {
		res.FuncEvaluations++
		return objective(x)
	}

	res.BestParticles = make([][]float64, 0, settings.MaxIterations)
	res.BestObjectives = make([]float64, 0, settings.MaxIterations)

	// initialize population with velocities and best known positions
	particles := make([][]float64, settings.PopulationSize)
	velocities := make([][]float64, settings.PopulationSize)
	bestPositions := make([][]float64, settings.PopulationSize)
	bestObjs := make([]float64, settings.PopulationSize)
	globalBest := make([]float64, 0)
	globalBestObj := math.MaxFloat64

	for i := range particles {

		particles[i], velocities[i] = init()
		bestObjs[i] = evaluate(particles[i])
		bestPositions[i] = make([]float64, len(particles[i]))
		copy(bestPositions[i], particles[i])

		if bestObjs[i] < globalBestObj {
			globalBest = make([]float64, len(particles[i]))
			copy(globalBest, particles[i])
			globalBestObj = bestObjs[i]
		}
	}

	res.BestObjectives = append(res.BestObjectives, globalBestObj)
	res.BestParticles = append(res.BestParticles, globalBest)

	for i := 0; i < settings.MaxIterations; i++ {
		totalObj := 0.0
		newGlobalBest := false
		newGlobalBestParticle := make([]float64, len(globalBest))
		for j, particle := range particles {
			velocity := velocities[j]
			for d, v := range velocity {
				rp, rg := rand.Float64(), rand.Float64()
				w := settings.Omega
				phip, phig := settings.ParticleWeight, settings.GlobalWeight
				velocity[d] = w*v + phip*rp*(bestPositions[j][d]-particle[d]) + phig*rg*(globalBest[d]-particle[d])
			}
			for d, p := range particle {
				particle[d] = p + settings.LearningRate*velocity[d]
			}
			obj := evaluate(particle)
			if obj < bestObjs[j] {
				copy(bestPositions[j], particle)
				bestObjs[j] = obj
				if obj < globalBestObj {
					newGlobalBest = true
					copy(newGlobalBestParticle, particle)
					copy(globalBest, particle)
					globalBestObj = obj
				}
			}
			totalObj += obj
		}
		if newGlobalBest {
			next := make([]float64, len(globalBest))
			copy(next, globalBest)
			res.BestParticles = append(res.BestParticles, next)
			res.BestObjectives = append(res.BestObjectives, globalBestObj)
		}
		logger.AddLine(i, []string{
			fmt.Sprint(i),
			fmt.Sprint(totalObj / float64(settings.PopulationSize)),
			fmt.Sprint(globalBestObj),
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
