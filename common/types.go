package common

import "time"

// Settings represents the settings of the optimization run
type Settings struct {
	// MaxIterations is the maximum number of iterations run by the algorithm
	// the algorithm will stop after this number is reached
	MaxIterations int
	// Verbose controls wether the algorithm should log information into the
	// console. 0 means no logging, n will log every n iterations
	Verbose int
}

// Result represents result information of the optimization, including
// statistics about the run
type Result struct {
	Runtime         time.Duration
	FuncEvaluations int
	Iterations      int
}
