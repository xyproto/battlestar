package battlestarlib

import (
	"strconv"
)

type (
	// Program is the source code that is to be compiled
	Program string

	// ProgramState is the state of the current position in this program, when compiling
	ProgramState struct {
		surpriseEndingWithExit bool           // To keep track of function blocks that are ended with "exit"
		loopStep               int            // To keep track of if rep should use stosb or stosw (and stepsize in loops in general)
		loopNameCounter        int            // To keep track of which generated label names have already been used
		ifNameCounter          int            // To keep track of which generated label names have already been used
		definedNames           []string       // all defined variables/constants/functions
		variables              map[string]int // map of variable names and reserved bytes
		inFunction             string         // name of the function we are currently in
		inLoop                 string         // name of the loop we are currently in
		inIfBlock              string         // name of the if block we are currently in
		endless                bool           // ending the program with endless keyword?
	}
)

var (
	// Global program state
	dataNotValueTypes []string          // all defined constants that are data (x: db 1,2,3,4...)
	types             map[string]string // type of the defined names
)

const (
	// For the types of loops that does not save and restore the counter before and after the loop body
	rawloopPrefix = "r_"
	// For the types of loops that loop forever
	endlessloopPrefix = "e_"
)

// NewProgramState returns a new state struct that is used when the program is compiled
func NewProgramState() *ProgramState {
	var ps ProgramState
	// Initialize global maps and slices
	ps.definedNames = make([]string, 0)
	ps.variables = make(map[string]int)
	return &ps
}

func (p *ProgramState) newLoopLabel() string {
	p.loopNameCounter++
	return "l" + strconv.Itoa(p.loopNameCounter)
}

func (p *ProgramState) newIfLabel() string {
	p.ifNameCounter++
	return "if" + strconv.Itoa(p.ifNameCounter)
}
