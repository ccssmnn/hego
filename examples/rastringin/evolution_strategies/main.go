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

	x0 := []float64{rand.Float64()*10.0 - 5.0, rand.Float64()*10.0 - 5.0}

	settings := hego.ESSettings{}
	settings.MaxIterations = 1000
	settings.Verbose = settings.MaxIterations / 10
	settings.NoiseSigma = 1.0
	settings.PopulationSize = 1000
	settings.LearningRate = 0.1

	result, err := hego.ES(rastringin, x0, settings)
	if err != nil {
		fmt.Printf("Got error while running Evolution Strategies Algorithm: %v", err)
	}
	fmt.Printf("Finished Evolution Strategies Algorithm! Result: %v, Value: %v \n", result.BestCandidate, result.BestObjective)
}
