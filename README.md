[![Gitpod ready-to-code](https://img.shields.io/badge/Gitpod-ready--to--code-blue?logo=gitpod)](https://gitpod.io/#https://github.com/ccssmnn/hego)

# hego

hego aims to be an efficient, zero dependency library for metaheuristic algorithms written in Go.

Finding the right algorithm and parameters for your problem is hard enough, you should not also bother implementing all these algorithms. hego aims to provide a rich set of helper functions to parametrize the algorithms for your needs as well as allowing you to provide your own functions e.g. for state changes.

## How?

Hego is able to solve your optimization problems when the algorithm specific interface is implemented. This approach allows you to use hego for various problem encodings. For example integer- or boolean vectors, graphs, structs etc.

For basic vector types (int, bool and float64) helper methods implemented in the subpackages `hego/crossover` and `hego/mutate` allow you to try out different parameter variants of the algorithms.

## Algorithms

- Simulated Annealing
- Genetic Algorithm

TODO:

- Ant Colony Optimization
- Glowworm Swarm Optimization

## Usage

The following example implements the `GeneticGenome` interface for the [Rastringin function](https://en.wikipedia.org/wiki/Rastrigin_function).

```golang
package main

import (
	"fmt"
	"math"
	"math/rand"

    "github.com/ccssmnn/hego"
	"github.com/ccssmnn/hego/mutate"
	"github.com/ccssmnn/hego/crossover"
)

func rastringin(x, y float64) float64 {
	return 10*2 + (x*x - 10*math.Cos(2*math.Pi*x)) + (y*y - 10*math.Cos(2*math.Pi*y))
}

type genome struct {
	v []float64
}

func (g *genome) GetGene() []interface{} {
	gene := make([]interface{}, len(g.v))
	for i, value := range g.v {
		gene[i] = value
	}
	return gene
}

func (g *genome) Crossover(other hego.GeneticGenome) hego.GeneticGenome {
	clone := genome{v: make([]float64, len(g.v))}
	gene := hego.ConvertFloat64(other.GetGene())
	clone.v = crossover.Arithmetic(g.v, gene, [2]float64{-0.5, 1.5})
	return &clone
}

func (g *genome) Mutate() hego.GeneticGenome {
	n := genome{v: make([]float64, len(g.v))}
	n.v = mutate.Gauss(g.v, 1.0)
	return &n
}

func (g *genome) Fitness() float64 {
	return rastringin(g.v[0], g.v[1])
}

func main() {
	// initialize population
	population := make([]hego.GeneticGenome, 100)
	for i := range population {
		population[i] = &genome{v: []float64{-10.0 + 10.0*rand.Float64(), -10 + 10*rand.Float64()}}
	}

	// set algorithm behaviour here
	settings := hego.GeneticSettings{}
	settings.MutationRate = 0.3
	settings.Elitism = 5
	settings.MaxIterations = 100
	settings.Verbose = 10

	result, err := hego.Genetic(population, settings)

	if err != nil {
		fmt.Printf("Got error while running Anneal: %v", err)
	} else {
		// extract solution and print result
		solution := result.BestGenome[result.Iterations-1].GetGene()
		value := result.BestFitness[result.Iterations-1]
		fmt.Printf("Finished Genetic Algorithm in %v! Needed %v function evaluations\n", result.Runtime, result.FuncEvaluations)
		fmt.Printf("Minimum found at x = [%v, %v] with f(x) = %v\n", solution[0], solution[1], value)
	}
	return
}
```

Logs:

```
Starting Genetic Algorithm...
   Iteration      Average Fitness           Best Fitness
           0    88.97123620542999     16.462962195504538
          10   21.401663002582232     1.0547108885271737
          20    24.59093664229806      0.877921289491681
          30    22.58149327408406     0.7019604847649106
          40   22.511145692890995    0.03094137018226739
          50    21.41978975257026    0.03094137018226739
          60    20.30242725390599   0.011517567887311841
          70   19.566038141815437   0.002379945277670714
          80   21.959566896637966   0.002379945277670714
          90     22.9504809485425   0.002379945277670714
          99   23.678249276698597   0.002379945277670714
DONE after 19.64255ms
Finished Genetic Algorithm in 19.64255ms! Needed 9361 function evaluations
Minimum found at x = [-0.003463162201316693, -5.6113457259983346e-05] with f(x) = 0.002379945277670714
```

For other usage examples checkout the examples directory. Examples are ordered by problem / algorithm. Currently implemented examples are:

- rastringin function, represantative continuous problem
- traveling salesman, permutation problem
- knapsack, binary selection problem

## Caveats

The Genetic Algorithm needs a type conversion / assertion from your custom genome to an interface slice and vice versa. I couldn't find another way to work with an interface while being able to handle any problem encoding. This has bad performance implications, but for convenience hego provides helper methods like `hego.ConvertFloat64`.
