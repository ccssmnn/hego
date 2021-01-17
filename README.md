[![Gitpod ready-to-code](https://img.shields.io/badge/Gitpod-ready--to--code-blue?logo=gitpod)](https://gitpod.io/#https://github.com/ccssmnn/hego) [![stability-unstable](https://img.shields.io/badge/stability-unstable-yellow.svg)](https://github.com/emersion/stability-badges#unstable) [![codecov](https://codecov.io/gh/ccssmnn/hego/branch/master/graph/badge.svg?token=F52EAPT69U)](https://codecov.io/gh/ccssmnn/hego) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

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

Hego is able to solve your optimization problems when the algorithm specific interface is implemented. This approach allows you to use hego for various problem encodings. For example integer- or boolean vectors, graphs, structs etc.

For basic vector types (int, bool and float64) helper methods implemented in the subpackages `hego/crossover` and `hego/mutate` allow you to try out different parameter variants of the algorithms.

Some algorithms however are only designed for specific sets of optimization problems. In these cases the algorithms provide an easier call signature that only requires the objective and the initial guess or initializer functions.

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

### Evolution Strategies

For [Evolution Strategies](https://openai.com/blog/evolution-strategies/) you only need to provide an objective function `func(v []float64) float64` and an initial vector `var x0 []float64`. Thats because Evolution Strategies is a gradient estimation strategy and therefore this algorithm only supports minimizing continuous problems.

```golang
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
	fmt.Printf("Finished Evolution Strategies Algorithm! Result: %v, Value: %v \n", result.Candidates[result.Iterations-1], result.BestObjective[result.Iterations-1])
	return
}

```

Logs:

```
Starting Evolution Strategy Algorithm...
   Iteration      Population Mean     Current Candidate
           0   42.224651908748704     6.368128212514245
         100   22.165302967018864    0.8194490791272049
         200   21.978639608437753   0.08232895412879238
         300   21.846276137579434   0.07416924829505689
         400    21.83424432866074    0.6522660494696222
         500   22.229350277127352   0.23431951698029962
         600   21.975966283045697     0.134608450035417
         700   21.774740603320893    0.5584185645513227
         800   21.754124247952106     0.594901719915752
         900   21.772541504660225    0.2556416152500045
         999   22.232498781270913    0.4453563633436506

DONE after 142.324509ms
Finished Evolution Strategies Algorithm! Result: [-0.016051413222810604 -0.0070487507112146994], Value: 0.4453563633436506
```

### Particle Swarm Optimization
For [Particle Swarm Optimization](https://en.wikipedia.org/wiki/Particle_swarm_optimization) you need to provide the objective function, one init method that initializes a particle and its velocity and the algorithm settings. PSO is also designed for continuous optimization problems.

```golang
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
	fmt.Printf("Finished Particle Swarm Optimization! Result: %v, Value: %v \n", result.BestParticles[len(result.BestParticles)-1], result.BestObjectives[len(result.BestParticles)-1])
	return
}
```
Logs
```
Iteration      Population Mean          Population Best
           0    37.02758959218644       1.0903639395271298
         100    12.17235438829838    3.296669603969349e-10
         200   11.984822426474523   2.8812507935072063e-12
         300    11.98986211949263   1.9966250874858815e-12
         400    11.16911682958532    8.331113576787175e-13
         500   11.324554904448338   2.0250467969162855e-13
         600   10.603598197357645    3.552713678800501e-15
         700   10.657698474188042                        0
         800   10.424831191887263                        0
         900   10.175484773908693                        0
         999   10.112701340657786                        0

Done after 13.589655ms!
Finished Particle Swarm Optimization! Result: [-1.6137381095924141e-09 -1.632641297328185e-09], Value: 0 
```

### Tabu Search

For [Tabu Search](https://en.wikipedia.org/wiki/Tabu_search) you need to implement the `TabuState` interface.

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

// state is a two element vector, it will implement the State interface for Tabu Search
type state []float64

// Neighbor produces another state by adding gaussian noise to the current state
func (s state) Neighbor() hego.TabuState {
	return state(mutate.Gauss(s, 0.3))
}

// Equal returns true, of two states are almost equal
func (s state) Equal(other hego.TabuState) bool {
	otherState := other.(state)
	for i := range s {
		if math.Abs(s[i]-otherState[i]) > 1e-4 {
			return false
		}
	}
	return true
}

// Objective of the current state. Lower is better
func (s state) Objective() float64 {
	return rastringin(s[0], s[1])
}

func main() {

	initialState := state{5.0, 5.0}

	settings := hego.TSSettings{}
	settings.MaxIterations = 1000
	settings.Verbose = settings.MaxIterations / 10
	settings.TabuListSize = 50
	settings.NeighborhoodSize = 25

	// start tabu search algorithm
	result, err := hego.TS(initialState, settings)

	if err != nil {
		fmt.Printf("Got error while running Anneal: %v", err)
	}
	finalState := result.States[len(result.States)-1]
	finalEnergy := finalState.Objective()
	fmt.Printf("Finished Tabu Search in %v! Result: %v, Value: %v \n", result.Runtime, finalState, finalEnergy)
	return
}
```

Logs:

```
   Iteration             Objective                   Best
           0     51.73729249013981      51.73729249013981
         100    10.928390644780535     10.449066694314698
         200     4.979290377705242     0.9958749953531694
         300    2.0591519405806906   0.052269383345688425
         400   0.21683395217992185   0.008851428787803428
         500     5.147812579878252   0.008851428787803428
         600    2.2625087237816253   0.000735083029045569
         700    1.3260815465805376   0.000735083029045569
         800      2.17002190874649   0.000735083029045569
         900    4.6177943136609745   0.000735083029045569
         999    0.2010011213275682   0.000735083029045569

Done after 7.11896ms!
Finished Tabu Search in 7.11896ms! Result: [-0.001324101190205762 0.0013971334657496629], Value: 0.000735083029045569
```

## Examples

hego contains a detailed examples directory. It is ordered by problem type and shows how to use hego to find solutions for these types of problems:

- Traveling Salesman Problem, an integer encoded permutation problem for finding the shortest path to visit all cities ([wikipedia](https://en.wikipedia.org/wiki/Travelling_salesman_problem))
- Knapsack Problem, a binary encoded combinatorical optimization problem to select items to get be best value while respecting the maximum weight ([wikipedia](https://en.wikipedia.org/wiki/Knapsack_problem))
- Rastrigin Function, a non convex function with a large number of local minima ([wikipedia](https://en.wikipedia.org/wiki/Rastrigin_function))

## Contributing

This repo is accepting PR's and welcoming issues. Feel free to contribute in any kind if

- you find any bugs
- you have ideas to make this library easier to use
- you have ideas on how to improve the performance
- you miss algorithm XY

Just be nice.

## License

The MIT License (MIT). [License](https://github.com/ccssmnn/hego)
