package main

import (
	"fmt"
	"math/rand"

	"github.com/ccssmnn/hego"
	"github.com/ccssmnn/hego/crossover"
	"github.com/ccssmnn/hego/mutate"
)

var nurses int
var shifts int
var days int
var maxNightShift int
var nightIndex int
var nursesPerShift []int
var shiftRequests [][][]bool // [d, s, n], true when n wants shift s on day d
var offRequests [][]bool     // [d, n], true when n wants off on day d

var shiftWeight float64
var offWeight float64

// scheduleQuality returns a quality value for respecting shift requests and off
// requests
func scheduleQuality(schedule [][][]bool) float64 {
	quality := 0.0
	for d := 0; d < days; d++ {
		for s := 0; s < shifts; s++ {
			for n := 0; n < nurses; n++ {
				if schedule[d][s][n] && shiftRequests[d][s][n] {
					quality += shiftWeight
				}
				if schedule[d][s][n] && offRequests[d][n] {
					quality -= offWeight
				}
			}
		}
	}
	return quality
}

var twoShiftsPenalty float64
var noNightShiftBreakPenalty float64

// schedulePenalty returns a penalty value for schedule properties that we want
// to avoid. e.g. having two shifts on the same day or working on a day after
// having a night shift
func schedulePenalty(schedule [][][]bool) float64 {
	penalty := 0.0
	for d := 0; d < days; d++ {
		shiftToday := make([]bool, nurses)
		for s := 0; s < shifts; s++ {
			for n := 0; n < nurses; n++ {
				if schedule[d][s][n] {
					// penalize two shifts on the same day
					if shiftToday[n] {
						penalty += twoShiftsPenalty
					} else {
						shiftToday[n] = true
					}
					// penalize having a day shift after a night shift
					if d > 0 && schedule[d-1][nightIndex][n] && s != nightIndex {
						penalty += noNightShiftBreakPenalty
					}
				}
			}
		}
	}
	return penalty
}

// we use a global schedule array for performance reasons
var schedule [][][]bool

func initSchedule() {
	schedule = make([][][]bool, days)
	for d := 0; d < days; d++ {
		schedule[d] = make([][]bool, shifts)
		for s := 0; s < shifts; s++ {
			schedule[d][s] = make([]bool, nurses)
		}
	}
}

// sets all assignments in this schedule to false
func resetSchedule() {
	for d := range schedule {
		for s := range schedule[d] {
			for n := range schedule[d][s] {
				schedule[d][s][n] = false
			}
		}
	}
}

type shift struct {
	day   int
	shift int
}

var allShifts []shift

// the genome is a list of ints. The index encodes the shift, the int is the
// nurse index for this shift
type genome []int

// Mutate either swaps two assignments or flips one assignment to another nurse
func (g genome) Mutate() hego.Genome {
	if rand.Float64() > 0.5 {
		return genome(mutate.Swap(g))
	}
	mutated := make(genome, len(g))
	mutated[rand.Intn(len(mutated))] = rand.Intn(nurses)
	return mutated
}

func (g genome) Crossover(other hego.Genome) hego.Genome {
	return genome(crossover.TwoPointInt(g, other.(genome)))
}

// resets schedule and writes nurse assignments into the schedule
func (g genome) fillSchedule() {
	resetSchedule()
	for i, assignment := range g {
		s := allShifts[i]
		schedule[s.day][s.shift][assignment] = true
	}
}

// Fitness is the quality and penalty for this schedules
func (g genome) Fitness() float64 {
	g.fillSchedule()
	return -scheduleQuality(schedule) + schedulePenalty(schedule)
}

func main() {
	// problem equivalent to https://developers.google.com/optimization/scheduling/employee_scheduling#requests
	nurses = 5
	shifts = 3
	days = 7
	nightIndex = 2
	maxNightShift = 7 // no max night shift
	nursesPerShift = []int{1, 1, 1}
	shiftWeight = 1.0              // weight for respecting a shift request
	offWeight = 2.0                // weight for respecting an off request
	twoShiftsPenalty = 5.0         // penalty for having two shifts on the same day
	noNightShiftBreakPenalty = 5.0 // penalty for working on a day after a night shift

	shiftRequests = make([][][]bool, days)
	offRequests = make([][]bool, days)
	for d := 0; d < days; d++ {
		shiftRequests[d] = make([][]bool, shifts)
		offRequests[d] = make([]bool, nurses)
		for s := 0; s < shifts; s++ {
			shiftRequests[d][s] = make([]bool, nurses)
		}
	}

	// nurse 0
	shiftRequests[0][0][0] = true
	shiftRequests[4][2][0] = true
	shiftRequests[5][1][0] = true
	shiftRequests[6][2][0] = true
	// nurse 1
	shiftRequests[2][1][1] = true
	shiftRequests[3][1][1] = true
	shiftRequests[4][0][1] = true
	shiftRequests[6][2][1] = true
	// nurse 2
	shiftRequests[0][1][2] = true
	shiftRequests[1][1][2] = true
	shiftRequests[3][0][2] = true
	shiftRequests[5][1][2] = true
	// nurse 3
	shiftRequests[0][2][3] = true
	shiftRequests[2][0][3] = true
	shiftRequests[3][1][3] = true
	shiftRequests[5][0][3] = true
	// nurse 4
	shiftRequests[1][2][4] = true
	shiftRequests[2][1][4] = true
	shiftRequests[4][0][4] = true
	shiftRequests[5][1][4] = true

	// init schedule
	initSchedule()

	// init shift slice containing all shifts that need to be assigned
	allShifts = make([]shift, 0)
	for d := 0; d < days; d++ {
		for s := 0; s < shifts; s++ {
			for n := 0; n < nursesPerShift[s]; n++ {
				allShifts = append(allShifts, shift{day: d, shift: s})
			}
		}
	}

	populationSize := 100
	settings := hego.GASettings{}
	settings.MaxIterations = 1000
	settings.Verbose = settings.MaxIterations / 10
	settings.MutationRate = 0.8
	settings.Elitism = 5
	settings.Selection = hego.RankBasedSelection

	initialPopulation := make([]hego.Genome, populationSize)

	for i := range initialPopulation {
		newGenome := make(genome, len(allShifts))
		for j := range allShifts {
			newGenome[j] = rand.Intn(nurses)
		}
		initialPopulation[i] = newGenome
	}

	// start genetic algorithm
	res, err := hego.GA(initialPopulation, settings)

	if err != nil {
		fmt.Printf("Got error while running Genetic Algorithm: %v", err)
		return
	}
	// extract result
	solution := res.BestGenome[len(res.BestGenome)-1].(genome)
	solution.fillSchedule()
	fmt.Printf("The solution found has an objective of %v \n", solution.Fitness())
	// print schedule
	for d := 0; d < days; d++ {
		fmt.Printf("\n\nDay %v", d)
		for s := 0; s < shifts; s++ {
			for n := 0; n < nurses; n++ {
				if schedule[d][s][n] && shiftRequests[d][s][n] {
					fmt.Printf("\nNurse %v works shift %v (requested)", n, s)
				} else if schedule[d][s][n] {
					fmt.Printf("\nNurse %v works shift %v (not requested)", n, s)
				}
			}
		}
	}
	fmt.Print("\n")
	return
}
