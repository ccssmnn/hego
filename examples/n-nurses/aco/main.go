package main

import (
	"fmt"
	"math"

	"github.com/ccssmnn/hego"
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
var pheromones = [][]float64{} // probability for each shift to be assigned to one nurse

type ant struct {
	index       int   // index of current shift
	assignments []int // which nurse is doing which shift
}

// Init resets internal state of ant
func (a *ant) Init() {
	a.index = 0
	a.assignments = make([]int, len(allShifts))
}

// Step adds `next` nurse to current day and current shift. When shift is full, goes to next shift/day. Returns true, when schedule is finished
func (a *ant) Step(next int) bool {
	a.assignments[a.index] = next
	a.index++
	return a.index == len(allShifts)
}

// PerceivePheromone returns pheromone values. For unfeasible selections the pheromone value is set to 0
func (a *ant) PerceivePheromone() []float64 {
	res := make([]float64, nurses)
	copy(res, pheromones[a.index][:])
	return res
}

// DropPheromone increases pheromone amount by 0.2 not considering the performance
func (a *ant) DropPheromone(performance float64) {
	for i, asmnt := range a.assignments {
		pheromones[i][asmnt] += 0.1
	}
}

// Evaporate increases pheromone amount by 0.2 not considering the performance
func (a *ant) Evaporate(factor, min float64) {
	for s := range allShifts {
		for n := 0; n < nurses; n++ {
			pheromones[s][n] = math.Max(min, pheromones[s][n]*factor)
		}
	}
}

// resets schedule and writes nurse assignments into the schedule
func (a *ant) fillSchedule() {
	resetSchedule()
	for i, assignment := range a.assignments {
		s := allShifts[i]
		schedule[s.day][s.shift][assignment] = true
	}
}

// Performance returns negative schedule quality
func (a *ant) Performance() float64 {
	a.fillSchedule()
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
	shiftRequests = make([][][]bool, days)
	offRequests = make([][]bool, days)
	shiftWeight = 1.0              // weight for respecting a shift request
	offWeight = 2.0                // weight for respecting an off request
	twoShiftsPenalty = 5.0         // penalty for having two shifts on the same day
	noNightShiftBreakPenalty = 5.0 // penalty for working on a day after a night shift

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

	initialPheromone := 1.0
	// reset pheromones
	pheromones = make([][]float64, len(allShifts))
	for s := range allShifts {
		p := make([]float64, nurses)
		for n := 0; n < nurses; n++ {
			p[n] = initialPheromone
		}
		pheromones[s] = p
	}

	population := make([]hego.Ant, 100)
	for i := range population {
		population[i] = &ant{}
	}

	settings := hego.ACOSettings{}
	settings.Evaporation = 0.99
	settings.MinPheromone = 0.1
	settings.MaxIterations = 1000
	settings.Verbose = settings.MaxIterations / 10 // log 10 steps to look at convergence behaviour
	result, err := hego.ACO(population, settings)

	if err != nil {
		fmt.Printf("Got error while running Ant Colony Optimization: %v", err)
		return
	}
	fmt.Printf("Finished Ant Colony Optimization in %v!", result.Runtime)

	result.BestAnts[len(result.BestAnts)-1].(*ant).fillSchedule()

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
