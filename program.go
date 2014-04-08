package main

import (
	"strconv"
)

type (
	Program string

	ProgramState struct {
		surprise_ending_with_exit bool                // To keep track of function blocks that are ended with "exit"
		loop_step                 int                 // To keep track of if rep should use stosb or stosw (and stepsize in loops in general)
		loop_name_counter         int                 // To keep track of which generated label names have already been used
		defined_names             []string            // all defined variables/constants/functions
		variables                 map[string][]string // list of variable names per function name
		in_function               string              // name of the function we are currently in
		in_loop                   string              // name of the loop we are currently in
	}
)

var (
	// Global program state
	data_not_value_types []string          // all defined constants that are data (x: db 1,2,3,4...)
	types                map[string]string // type of the defined names
)

const (
	// For the types of loops that does not save and restore the counter before and after the loop body
	rawloop_prefix = "r_"
	// For the types of loops that loop forever
	endlessloop_prefix = "e_"
)

func NewProgramState() *ProgramState {
	var ps ProgramState
	// Initialize global maps and slices
	ps.defined_names = make([]string, 0, 0)
	ps.variables = make(map[string][]string)
	return &ps
}

func (p *ProgramState) new_loop_label() string {
	p.loop_name_counter += 1
	return "l" + strconv.Itoa(p.loop_name_counter)
}
