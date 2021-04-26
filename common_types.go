// Package hego provides methods to blackbox optimization algorithms like
// Genetic Algorithm, Simulated Annealing or Ant Colony Optimization
//
// The consistent API between algorithms and (optionally) verbose execution
// make finding the best algorithm and the right parameters easy and quick
package hego

import "time"

// Settings represents the settings of the optimization run
type Settings struct {
	// MaxIterations is the maximum number of iterations run by the algorithm
	// the algorithm will stop after this number is reached
	MaxIterations int
	// Verbose controls wether the algorithm should log information into the
	// console. 0 means no logging, n will log every n iterations
	Verbose int
	// KeepHistory, when true intermediate results are stored
	KeepHistory bool
}

// Result represents result information of the optimization, including
// statistics about the run
type Result struct {
	Runtime         time.Duration
	FuncEvaluations int
	Iterations      int
}
