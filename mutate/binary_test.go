package mutate

import "testing"

func TestFlip(t *testing.T) {
	state := []bool{false, true, true, false, true}
	neighbor := Flip(state)
	count := 0
	for i := 0; i < len(state); i++ {
		if state[i] != neighbor[i] {
			count++
		}
	}
	if count != 1 {
		t.Errorf("Expected exactly one bit to be flipped, found %v", count)
	}
}

func TestFlipn(t *testing.T) {
	state := []bool{false, true, true, false, true}
	neighbor := Flipn(state, 2)
	count := 0
	for i := 0; i < len(state); i++ {
		if state[i] != neighbor[i] {
			count++
		}
	}
	if count != 2 {
		t.Errorf("Expected exactly one bit to be flipped, found %v", count)
	}
}
