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

const shiftWeight = 1.0
const offWeight = 2.0

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

var bestPerformance = math.MaxFloat64
var pheromones = [7][3][5]float64{}

type ant struct {
	day                int        // current day
	shift              int        // current shift
	nurseCount         int        // current nursecount for shift and day
	shiftPerNurse      []int      // number of shifts for each nurse
	nightShiftPerNurse []int      // number of night shifts for each nurse
	schedule           [][][]bool // solution
}

// Init resets current state of ant
func (a *ant) Init() {
	a.day = 0
	a.shift = 0
	a.nurseCount = 0
	a.shiftPerNurse = make([]int, nurses)
	a.nightShiftPerNurse = make([]int, nurses)
	a.schedule = make([][][]bool, days)
	for d := 0; d < days; d++ {
		shift := make([][]bool, shifts)
		for s := 0; s < shifts; s++ {
			assignment := make([]bool, nurses)
			shift[s] = assignment
		}
		a.schedule[d] = shift
	}
}

// Step adds `next` nurse to current day and current shift. When shift is full, goes to next shift/day. Returns true, when schedule is finished
func (a *ant) Step(next int) bool {
	// reset when next is invalid
	if next == -1 {
		a.Init()
		return false
	}
	a.schedule[a.day][a.shift][next] = true
	a.shiftPerNurse[next]++
	if a.shift == nightIndex {
		a.nightShiftPerNurse[next]++
	}
	a.nurseCount++
	// go to next shift
	if a.nurseCount == nursesPerShift[a.shift] {
		a.nurseCount = 0
		a.shift++
	}
	// go to next day
	if a.shift == shifts {
		a.shift = 0
		a.day++
	}
	return a.day == days
}

// PerceivePheromone returns pheromone values. For unfeasible selections the pheromone value is set to 0
func (a *ant) PerceivePheromone() []float64 {
	res := make([]float64, nurses)
	copy(res, pheromones[a.day][a.shift][:])
	// increase probability if nurse has requested this shift
	// and reduce if nurse has requested day off
	for s := 0; s < shifts; s++ {
		for n := 0; n < nurses; n++ {
			if shiftRequests[a.day][s][n] {
				res[n] *= 2.0
			}
			if offRequests[a.day][n] {
				res[n] /= 2.0
			}
		}
	}
	// do not take nurse, that already has a shift this day
	for s := 0; s < shifts; s++ {
		for n := 0; n < nurses; n++ {
			if a.schedule[a.day][s][n] {
				res[n] = 0.0
			}
		}
	}
	// do not take nurse, that had night shift the day before or has reached maxNightShift
	if a.day > 0 {
		for n := 0; n < nurses; n++ {
			if a.schedule[a.day-1][nightIndex][n] || a.nightShiftPerNurse[n] == maxNightShift {
				res[n] = 0.0
			}
		}
	}
	return res
}

// DropPheromone increases pheromone amount by 0.2 not considering the performance
func (a *ant) DropPheromone(performance float64) {
	for d := 0; d < days; d++ {
		for s := 0; s < shifts; s++ {
			for n := 0; n < nurses; n++ {
				pheromones[d][s][n] = 1 / (1 + (bestPerformance-performance)/bestPerformance)
			}
		}
	}
}

// Evaporate increases pheromone amount by 0.2 not considering the performance
func (a *ant) Evaporate(factor, min float64) {
	for d := 0; d < days; d++ {
		for s := 0; s < shifts; s++ {
			for n := 0; n < nurses; n++ {
				pheromones[d][s][n] = math.Max(min, pheromones[d][s][n]*factor)
			}
		}
	}
}

// Performance returns negative schedule quality
func (a *ant) Performance() float64 {
	performance := -scheduleQuality(a.schedule)
	if performance < bestPerformance {
		bestPerformance = performance
	}
	return performance
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

	initialPheromone := 1.0
	// reset pheromones
	for d := 0; d < days; d++ {
		for s := 0; s < shifts; s++ {
			for n := 0; n < nurses; n++ {
				pheromones[d][s][n] = initialPheromone
			}
		}
	}
	population := make([]hego.Ant, 100)
	for i := range population {
		population[i] = &ant{}
	}
	settings := hego.ACOSettings{}
	settings.Evaporation = 0.999
	settings.MinPheromone = 0.01
	settings.MaxIterations = 10000
	settings.Verbose = settings.MaxIterations / 10 // log 10 steps to look at convergence behaviour
	result, err := hego.ACO(population, settings)

	if err != nil {
		fmt.Printf("Got error while running Ant Colony Optimization: %v", err)
		return
	}
	fmt.Printf("Finished Ant Colony Optimization in %v! Needed %v function evaluations\n", result.Runtime, result.FuncEvaluations)

	finalSchedule := result.BestAnts[len(result.BestAnts)-1].(*ant).schedule
	// print schedule
	for d := 0; d < days; d++ {
		fmt.Printf("\n\nDay %v", d)
		for s := 0; s < shifts; s++ {
			for n := 0; n < nurses; n++ {
				if finalSchedule[d][s][n] && shiftRequests[d][s][n] {
					fmt.Printf("\nNurse %v works shift %v (requested)", n, s)
				} else if finalSchedule[d][s][n] {
					fmt.Printf("\nNurse %v works shift %v (not requested)", n, s)
				}
			}
		}
	}
	fmt.Print("\n")
	return
}
