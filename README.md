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

The following example implements the `Genome` interface for the [Rastringin function](https://en.wikipedia.org/wiki/Rastrigin_function).

```golang
package main

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/ccssmnn/hego"
	"github.com/ccssmnn/hego/crossover"
	"github.com/ccssmnn/hego/mutate"
)

func rastringin(x, y float64) float64 {
	return 10*2 + (x*x - 10*math.Cos(2*math.Pi*x)) + (y*y - 10*math.Cos(2*math.Pi*y))
}

// genome is a vector of float values
type genome []float64
<<<<<<< HEAD

// Crossover returns a child genome which is a combination of the current and other
// genome. Here an the Arithmetic crossover operation is used
func (g genome) Crossover(other hego.Genome) hego.Genome {
	return genome(crossover.Arithmetic(g, other.(genome), [2]float64{-0.5, 1.5}))
}

// Mutate adds variation to a genome. In this case we add gaussian noise
func (g genome) Mutate() hego.Genome {
	return genome(mutate.Gauss(g, 0.5))
}

=======

// Crossover returns a child genome which is a combination of the current and other
// genome. Here an the Arithmetic crossover operation is used
func (g genome) Crossover(other hego.Genome) hego.Genome {
	child := genome(crossover.Arithmetic(g, *other.(*genome), [2]float64{-0.5, 1.5}))
	return &child
}

// Mutate adds variation to a genome. In this case we add gaussian noise
func (g genome) Mutate() hego.Genome {
	mutant := genome(mutate.Gauss(g, 0.5))
	return &mutant
}

>>>>>>> b429d661041a7adbd345995f54c15f84c74acebb
// Fitness is called to evaluate the objective functino value. Lower is better
func (g genome) Fitness() float64 {
	return rastringin(g[0], g[1])
}

func main() {
	// initialize population
	population := make([]hego.Genome, 100)
	for i := range population {
<<<<<<< HEAD
		population[i] = genome{-10.0 + 10.0*rand.Float64(), -10 + 10*rand.Float64()}
=======
		population[i] = &genome{-10.0 + 10.0*rand.Float64(), -10 + 10*rand.Float64()}
>>>>>>> b429d661041a7adbd345995f54c15f84c74acebb
	}

	// set algorithm behaviour here
	settings := hego.GASettings{}
	settings.MutationRate = 0.5
	settings.Elitism = 5
	settings.MaxIterations = 100
	settings.Verbose = 10

	// call genetic algorithm
	result, err := hego.GA(population, settings)

	if err != nil {
		fmt.Printf("Got error while running Genetic Algorithm: %v", err)
		return
	}

	// extract solution and print result
<<<<<<< HEAD
	solution := result.BestGenome[result.Iterations-1].(genome)
=======
	solution := *result.BestGenome[result.Iterations-1].(*genome)
>>>>>>> b429d661041a7adbd345995f54c15f84c74acebb
	value := result.BestFitness[result.Iterations-1]
	fmt.Printf("Finished Genetic Algorithm in %v! Needed %v function evaluations\n", result.Runtime, result.FuncEvaluations)
	fmt.Printf("Minimum found at x = [%v, %v] with f(x) = %v\n", solution[0], solution[1], value)
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
