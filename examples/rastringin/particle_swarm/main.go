package main

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/ccssmnn/hego"
)

func rastringin(v []float64) float64 {
	x, y := v[0], v[1]
	return 10*2 + (x*x - 10*math.Cos(2*math.Pi*x)) + (y*y - 10*math.Cos(2*math.Pi*y))
}

func main() {

	init := func() ([]float64, []float64) {
		return []float64{
				rand.Float64()*10.0 - 5.0,
				rand.Float64()*10.0 - 5.0,
			}, []float64{
				rand.Float64(),
				rand.Float64(),
			}
	}

	settings := hego.PSOSettings{}
	settings.MaxIterations = 1000
	settings.Verbose = settings.MaxIterations / 10
	settings.PopulationSize = 100
	settings.GlobalWeight = 0.3
	settings.Omega = 0.5
	settings.ParticleWeight = 0.1
	settings.LearningRate = 0.5

	result, err := hego.PSO(rastringin, init, settings)
	if err != nil {
		fmt.Printf("Got error while running Particle Swarm Optimization: %v", err)
	}
	fmt.Printf("Finished Particle Swarm Optimization! Result: %v, Value: %v \n", result.BestParticle, result.BestObjective)
}
