[![Gitpod ready-to-code](https://img.shields.io/badge/Gitpod-ready--to--code-blue?logo=gitpod)](https://gitpod.io/#https://github.com/ccssmnn/hego) [![stability-unstable](https://img.shields.io/badge/stability-unstable-yellow.svg)](https://github.com/emersion/stability-badges#unstable) [![codecov](https://codecov.io/gh/ccssmnn/hego/branch/master/graph/badge.svg?token=F52EAPT69U)](https://codecov.io/gh/ccssmnn/hego) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT) [![Go Report Card](https://goreportcard.com/badge/github.com/ccssmnn/hego)](https://goreportcard.com/report/github.com/ccssmnn/hego) [![Go Reference](https://pkg.go.dev/badge/github.com/ccssmnn/hego.svg)](https://pkg.go.dev/github.com/ccssmnn/hego)

# hego

hego aims to provide a consistent API for several metaheuristics (black box optimization algorithms) while being performant.

Even though most of the metaheuristics might fit to any kind of optimization problem most of them have some caveats / advantages in different fields. hego allows you to rapidly try different algorithms and experiment with the parameters in order to solve your problems in the best possible time-to-quality ratio.

## Algorithms

Currently the following algorithms are implemented:

- Simulated Annealing (SA)
- Genetic Algorithm (GA)
- Ant Colony Optimization (ACO)
- Tabu Search (TS)
- Evolution Strategies (ES) (continuous only)
- Particle Swarm Optimization (PSO) (continuous only)

All algorithms are implemented for finding minimum values.

## Usage

hego is able to solve your optimization problems when the algorithm specific interface is implemented. This approach allows you to use hego for various problem encodings. For example integer- or boolean vectors, graphs, structs etc.

For basic vector types (int, bool and float64) helper methods implemented in the subpackages `hego/crossover` and `hego/mutate` allow you to experiment with different parameter variants of the algorithms.

Some algorithms however are only designed for specific sets of optimization problems. In these cases the algorithms provide an easier call signature that only requires the objective and the initial guess or initializer functions. (Evolution Strategies, Particle Swarm Optimization)

hego has a rich examples directory. It is ordered by problem type and shows how to apply hego's algorithms to these types of problems:

- Traveling Salesman Problem, an integer encoded permutation problem for finding the shortest path to visit all cities ([wikipedia](https://en.wikipedia.org/wiki/Travelling_salesman_problem))
- Knapsack Problem, a binary encoded combinatorical optimization problem to select items to get be best value while respecting the maximum weight ([wikipedia](https://en.wikipedia.org/wiki/Knapsack_problem))
- Rastrigin Function, a non convex function with a large number of local minima ([wikipedia](https://en.wikipedia.org/wiki/Rastrigin_function))
- Nurse Scheduling Problem, a scheduling problem for assigning shifts to nurses with constraints ([wikipedia](https://en.wikipedia.org/wiki/Nurse_scheduling_problem))
- Vehicle Routing Problem, a combination of Knapsack and Traveling Salesman problem ([wikipedia](https://en.wikipedia.org/wiki/Vehicle_routing_problem))

## Example

This example uses Simulated Annealing (SA) for optimizing the Rastrigin Function:

```golang
package main

import (
	"fmt"
	"math"

	"github.com/ccssmnn/hego"
	"github.com/ccssmnn/hego/mutate"
)

func rastringin(x, y float64) float64 {
	return 10*2 + (x*x - 10*math.Cos(2*math.Pi*x)) + (y*y - 10*math.Cos(2*math.Pi*y))
}

// state is a two element vector, it will implement the State interface for Simulated Annealing
type state []float64

// Neighbor produces another state by adding gaussian noise to the current state
func (s state) Neighbor() hego.AnnealingState {
	return state(mutate.Gauss(s, 0.3))
}

// Energy returns the energy of the current state. Lower is better
func (s state) Energy() float64 {
	return rastringin(s[0], s[1])
}

func main() {
	initialState := state{5.0, 5.0}

	settings := hego.SASettings{}
	settings.MaxIterations = 100000
	settings.Verbose = 10000
	settings.Temperature = 10.0       // choose temperature in the range of the systems energy
	settings.AnnealingFactor = 0.9999 // decrementing the temperature leads to convergence, we want to reach convergence when approaching the end of iterations

	// start simulated annealing algorithm
	result, err := hego.SA(initialState, settings)

	if err != nil {
		fmt.Printf("Got error while running Anneal: %v", err)
	}
	finalState := result.State
	finalEnergy := result.Energy
	fmt.Printf("Finished Simulated Annealing in %v! Result: %v, Value: %v \n", result.Runtime, finalState, finalEnergy)
}
```

It logs:

```
   Iteration             Temperature                  Energy
           0                   9.999                      50
       10000      3.6782426032832705      3.0986994133146712
       20000      1.3530821730781113       4.227542078387473
       30000       0.497746224098313       2.336322174938326
       40000      0.1831014468548652     0.30639618340376096
       50000     0.06735588984342127     0.03177535766328887
       60000    0.024777608121224735     0.02194743246350228
       70000    0.009114716851579779   0.0030078958948340784
       80000    0.003352949278962375    0.012710941747947402
       90000   0.0012334194303957732    0.004538678651899275
       99999   0.0004537723395901116   0.0008388313144696014

Done after 108.647155ms!
Finished Simulated Annealing in 108.647155ms! Result: [0.0010647353926910566 -0.001759125670646859], Value: 0.0008388313144696014
```

## Contributing

This repo is accepting PR's and welcoming issues. Feel free to contribute in any kind if

- you find any bugs
- you have ideas to make this library easier to use
- you have ideas on how to improve the performance
- you miss algorithm XY

Just be nice.

## License

The MIT License (MIT). [License](https://github.com/ccssmnn/hego)
