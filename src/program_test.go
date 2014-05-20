package main

import (
	"testing"
)

func TestProgramState(t *testing.T) {
	ps := NewProgramState()
	// Very general test, not too useful
	if len(ps.variables) != 0 {
		t.Errorf("Error initializing program state.\n")
	}
}

