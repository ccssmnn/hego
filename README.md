[![Gitpod ready-to-code](https://img.shields.io/badge/Gitpod-ready--to--code-blue?logo=gitpod)](https://gitpod.io/#https://github.com/ccssmnn/hego)

# hego

hego aims to provide a consistent API for several metaheuristics (black box optimization algorithms) while being performant.

Even though most of the metaheuristics might fit to any kind of optimization problem most of them have some caveats / advantages in different fields. hego allows you to rapidly try different algorithms and experiment with the parameters in order to solve your problems in the best possible time-to-quality ratio.

## Usage

Hego is able to solve your optimization problems when the algorithm specific interface is implemented. This approach allows you to use hego for various problem encodings. For example integer- or boolean vectors, graphs, structs etc.

For basic vector types (int, bool and float64) helper methods implemented in the subpackages `hego/crossover` and `hego/mutate` allow you to try out different parameter variants of the algorithms.

## Algorithms

Currently the following algorithms are implemented:

- Simulated Annealing (SA)
- Genetic Algorithm (GA)

These are in our scope (TODO):

- Ant Colony Optimization (ACO), good for permutation based problems
- Glowworm Swarm Optimization (GSO), nice for finding multiple local minima
- Evolutionary Strategies (ES), good for real-valued functions
- Memetic Algorithm (MA), Genetic Algorithm + Local Search

All algorithms are implement for finding minimum values.

The following examples show how to use hego to find solutions for the [Rastringin function](https://en.wikipedia.org/wiki/Rastrigin_function).

### Simulated Annealing

For [Simulated Annealing](https://en.wikipedia.org/wiki/Simulated_annealing) you need to implement the `AnnealingState` interface:

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
	fmt.Printf("Finished Simulated Annealing in %v! Result: %v, Value: %v \n", result.Runtime, result.States[result.Iterations], result.Energies[result.Iterations])
	return
}

```

Logs:

```
Starting Simulated Annealing...
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
DONE after 43.418236ms
Finished Simulated Annealing in 43.418236ms! Result: [0.0010647353926910566 -0.001759125670646859], Value: 0.0008388313144696014
```

### Genetic Algorithm

For the [Genetic Algorithm](https://en.wikipedia.org/wiki/Genetic_algorithm) you have to implement the `Genome` interface:

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

// Crossover returns a child genome which is a combination of the current and other
// genome. Here an the Arithmetic crossover operation is used
func (g genome) Crossover(other hego.Genome) hego.Genome {
	return genome(crossover.Arithmetic(g, other.(genome), [2]float64{-0.5, 1.5}))
}

// Mutate adds variation to a genome. In this case we add gaussian noise
func (g genome) Mutate() hego.Genome {
	return genome(mutate.Gauss(g, 0.5))
}

// Fitness is called to evaluate the objective functino value. Lower is better
func (g genome) Fitness() float64 {
	return rastringin(g[0], g[1])
}

func main() {
	// initialize population
	population := make([]hego.Genome, 100)
	for i := range population {
		population[i] = genome{-10.0 + 10.0*rand.Float64(), -10 + 10*rand.Float64()}
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
	solution := result.BestGenome[result.Iterations-1].(genome)
	fmt.Printf("Finished Genetic Algorithm in %v! Needed %v function evaluations\n", result.Runtime, result.FuncEvaluations)
	fmt.Printf("Minimum found at x = [%v, %v] with f(x) = %v\n", solution[0], solution[1], rastringin(solution[0], solution[1]))
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

## Examples

hego contains a detailed examples directory. It is ordered by problem type and shows how to use hego to find solutions for these types of problems:

- Traveling Salesman Problem, an integer encoded permutation problem for finding the shortest path to visit all cities ([wikipedia](https://en.wikipedia.org/wiki/Travelling_salesman_problem))
- Knapsack Problem, a binary encoded combinatorical optimization problem to select items to get be best value while respecting the maximum weight ([wikipedia](https://en.wikipedia.org/wiki/Knapsack_problem))
- Rastrigin Function, a non convex function with a large number of local minima ([wikipedia](https://en.wikipedia.org/wiki/Rastrigin_function))

## TODO

- Implement missing algorithms
- Find and implement more nature inspired heuristics
- Extend collection of mutation and crossover methods
- Add more real-world examples
  - constrained optimization
  - vehicle routing
  - multiple knapsack problem

## License

The MIT License (MIT). [License](https://github.com/ccssmnn/hego)
