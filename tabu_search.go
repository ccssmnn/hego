package hego

import (
	"fmt"
	"math"
	"time"
)

// TabuState describes the state during a tabu search
type TabuState interface {
	// Objective is the function to be minimized
	Objective() float64
	// Equal returns true, when other state is equal to current one
	Equal(other TabuState) bool
	// Neighbor produces a related state for local search
	Neighbor() TabuState
}

// TSResult holds result and progress information about the tabu search algorithm
type TSResult struct {
	// States holds the best states. Last element in this list is overall best solution
	States []TabuState
	// Objectives holds the best objectives. Each entry corresponds to an element in States
	Objectives    []float64
	BestState     TabuState
	BestObjective float64
	Result
}

// TSSettings describes the necessary settings for the tabu search algorithm
type TSSettings struct {
	// NeighborhoodSize sets the number of neighbors created in each iteration
	NeighborhoodSize int
	// TabuListSize is the memory of the algorithm. Each iteration the state
	// is added to the tabu list. A produced neighbor wont be selected if he appears
	// in the tabu list
	TabuListSize int
	Settings
}

// Verify returns an error if settings verification fails
func (s *TSSettings) Verify() error {
	if s.NeighborhoodSize <= 1 {
		return fmt.Errorf("size of neighborhood must be greater that 1, got %v", s.NeighborhoodSize)
	}
	if s.TabuListSize <= 1 {
		return fmt.Errorf("size of Tabu List must be larger than 1, got %v", s.TabuListSize)
	}
	return nil
}

// TS performs tabu search optimization
func TS(
	initialState TabuState,
	settings TSSettings,
) (res TSResult, err error) {

	err = settings.Verify()
	if err != nil {
		err = fmt.Errorf("settings verification failed: %v", err)
		return
	}

	start := time.Now()

	logger := newLogger("Tabu Search", []string{"Iteration", "Objective", "Best"}, settings.Verbose, settings.MaxIterations)

	evaluate := func(s TabuState) float64 {
		res.FuncEvaluations++
		return s.Objective()
	}

	state := initialState
	var obj float64
	tabuList := make([]TabuState, 0)

	inList := func(s TabuState) bool {
		for _, ts := range tabuList {
			if ts.Equal(s) {
				return true
			}
		}
		return false
	}

	if settings.KeepIntermediateResults {
		res.States = make([]TabuState, 0, settings.MaxIterations)
		res.Objectives = make([]float64, 0, settings.MaxIterations)
	}

	res.BestObjective = math.MaxFloat64

	for i := 0; i < settings.MaxIterations; i++ {

		bestNeighbor := state.Neighbor()
		bestNeighborObj := evaluate(bestNeighbor)

		for j := 0; j < settings.NeighborhoodSize; j++ {
			candidate := state.Neighbor()
			candidateObj := evaluate(candidate)
			if candidateObj < bestNeighborObj && !inList(candidate) {
				bestNeighbor = candidate
				bestNeighborObj = candidateObj
			}
		}

		tabuList = append(tabuList, bestNeighbor)

		if len(tabuList) > settings.TabuListSize {
			tabuList = tabuList[1:]
		}

		state = bestNeighbor
		obj = bestNeighborObj

		if settings.KeepIntermediateResults && (len(res.Objectives) == 0 || res.Objectives[len(res.Objectives)-1] > bestNeighborObj) {
			res.States = append(res.States, state)
			res.Objectives = append(res.Objectives, obj)
		}

		if res.BestObjective > obj {
			res.BestObjective = obj
			res.BestState = state
		}

		res.Iterations++
		if i == 0 {
			logger.AddLine(i, []string{
				fmt.Sprint(i),
				fmt.Sprint(obj),
				fmt.Sprint(obj),
			})
		} else {
			logger.AddLine(i, []string{
				fmt.Sprint(i),
				fmt.Sprint(obj),
				fmt.Sprint(res.BestObjective),
			})
		}
	}

	res.Runtime = time.Since(start)
	res.Iterations = settings.MaxIterations

	logger.Flush()
	if settings.Verbose > 0 {
		fmt.Printf("Done after %v!\n", res.Runtime)
	}
	return
}
