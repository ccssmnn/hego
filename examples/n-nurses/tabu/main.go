package main

import (
	"fmt"
	"math/rand"

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

// shiftInSlice helps checking if two states are equal
func shiftInSlice(s shift, slice []shift) bool {
	for i := range slice {
		if s.day == slice[i].day && s.shift == slice[i].shift {
			return true
		}
	}
	return false
}

// state is the tabu state that we use for optimization
// each nurse has one shift slice with shifts that this nurse is assigned to
type state struct {
	nurses [][]shift
}

func (s state) Equal(other hego.TabuState) bool {
	o := other.(state)
	for n := 0; n < nurses; n++ {
		if len(s.nurses[n]) != len(o.nurses[n]) {
			return false
		}
		for i := range s.nurses[n] {
			if !shiftInSlice(s.nurses[n][i], o.nurses[n]) {
				return false
			}
		}
	}
	return true
}

// swapShift chooses two nurses and swaps a shift between these nurses
func (s state) swapShift() hego.TabuState {
	var na, nb int
	for na == nb {
		na, nb = rand.Intn(len(s.nurses)), rand.Intn(len(s.nurses))
	}
	sha, shb := rand.Intn(len(s.nurses[na])), rand.Intn(len(s.nurses[nb]))
	neighbor := state{}
	neighbor.nurses = make([][]shift, len(s.nurses))
	for i := range neighbor.nurses {
		neighbor.nurses[i] = make([]shift, len(s.nurses[i]))
		copy(neighbor.nurses[i], s.nurses[i])
	}
	neighbor.nurses[na][sha] = s.nurses[nb][shb]
	neighbor.nurses[nb][shb] = s.nurses[na][sha]
	return neighbor
}

// transferShift transfers a shift from one nurse to another one
func (s state) transferShift() hego.TabuState {
	neighbor := state{}
	neighbor.nurses = make([][]shift, len(s.nurses))
	for i := range neighbor.nurses {
		neighbor.nurses[i] = make([]shift, len(s.nurses[i]))
		copy(neighbor.nurses[i], s.nurses[i])
	}
	var destination, source int
	for destination == source {
		destination, source = rand.Intn(len(s.nurses)), rand.Intn(len(s.nurses))
	}
	if len(s.nurses[destination]) > len(s.nurses[source]) {
		destination, source = source, destination
	}
	lenSource := len(neighbor.nurses[source])
	sh := rand.Intn(lenSource)
	neighbor.nurses[destination] = append(neighbor.nurses[destination], neighbor.nurses[source][sh])
	// deletes [sh] from source
	neighbor.nurses[source][sh] = neighbor.nurses[source][lenSource-1]
	neighbor.nurses[source] = neighbor.nurses[source][:lenSource-1]
	return neighbor
}

// Neighbor calls transfershift or swapshift
func (s state) Neighbor() hego.TabuState {
	if rand.Float64() < 0.5 {
		return s.swapShift()
	}
	return s.transferShift()
}

// resets schedule and writes nurse assignments into the schedule
func (s state) fillSchedule() {
	resetSchedule()
	for n, nShifts := range s.nurses {
		for _, s := range nShifts {
			schedule[s.day][s.shift][n] = true
		}
	}
}

// Objective is the quality and penalty for this schedules
func (s state) Objective() float64 {
	s.fillSchedule()
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

	// select initial state
	initialState := state{
		nurses: make([][]shift, nurses),
	}
	// init shift slice containing all shifts that need to be assigned
	allShifts := make([]shift, 0)
	for d := 0; d < days; d++ {
		for s := 0; s < shifts; s++ {
			for n := 0; n < nursesPerShift[s]; n++ {
				allShifts = append(allShifts, shift{day: d, shift: s})
			}
		}
	}
	// randomly assign shifts to nurses
	for _, s := range allShifts {
		n := rand.Intn(nurses)
		initialState.nurses[n] = append(initialState.nurses[n], s)
	}

	settings := hego.TSSettings{}
	settings.MaxIterations = 1000
	settings.Verbose = settings.MaxIterations / 10
	settings.TabuListSize = 10
	settings.NeighborhoodSize = 10

	// start tabu search main algorithm
	res, err := hego.TS(&initialState, settings)

	if err != nil {
		fmt.Printf("Got error while running Tabu Search: %v", err)
		return
	}
	// extract result
	solution := res.States[len(res.States)-1].(state)
	solution.fillSchedule()
	fmt.Printf("The solution found has an objective of %v \n", solution.Objective())
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
