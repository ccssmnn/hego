package mutate

import "testing"

func TestSwap(t *testing.T) {
	state := []int{0, 1, 2, 3, 4}
	neighbor := Swap(state)
	count := 0
	for i := 0; i < len(state); i++ {
		if state[i] != neighbor[i] {
			count++
		}
	}
	if count != 2 {
		t.Errorf("Swap should change exatcy two values, found %v", count)
	}
}

func TestSwapClose(t *testing.T) {
	state := []int{0, 1, 2, 3, 4}
	neighbor := SwapClose(state)
	indizes := make([]int, 0, 2)
	for i := 0; i < len(state); i++ {
		if state[i] != neighbor[i] {
			indizes = append(indizes, i)
		}
	}
	if len(indizes) != 2 {
		t.Errorf("SwapClose should change two neighbors, found %v", len(indizes))
	}
	if (indizes[0]+1)%len(state) != indizes[1] {
		t.Errorf("Swapped values should be neighbors. Got indizes %v and %v.", indizes[0], indizes[1])
	}
}
