package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

/*
	A,4=^2|B
	When in state A, and input is 4,
	Pop 2 from stack and goto state B.

	B,2=!4|A
	When in state B, and input is 2,
	Push 4 on stack, and goto state A.
*/

// (Current) state form
type StateForm struct {
	CurrentState string // The current state
	InputValue   int    // The value from input tape
}

// Transition form
type OpForm struct {
	Operation  string // The operation push/pop
	StackValue int    // To value to push, or pop
	NextState  string // Transitiotn to next state
}

type PdaMap map[StateForm]OpForm

const (
	TokComma   = ","
	TokEqual   = "="
	TokPush    = "!"
	TokPop     = "^"
	TokPipe    = "|"
	TokNewLine = "\n"
	TokEmpty   = ""
)

const ParseError = "parsing newline error!"

// Split node into two child nodes
func splitNode(str string, tok string) (string, string) {
	v := strings.Split(str, tok)
	if len(v) < 2 {
		return TokEmpty, TokEmpty
	}
	return v[0], v[1]
}

// Split the stack operation node
func splitStackop(str string) (string, string) {
	if strings.Contains(str, TokPush) {
		return TokPush, strings.TrimPrefix(str, TokPush)
	} else if strings.Contains(str, TokPop) {
		return TokPop, strings.TrimPrefix(str, TokPop)
	}
	return TokEmpty, TokEmpty
}

// Parse string, construct forms based on syntax tree
func parseLine(n0 string) (*StateForm, *OpForm, error) {
	var (
		sf  = new(StateForm)
		of  = new(OpForm)
		err error
	)

	// Line should not be empty
	if n0 == TokEmpty {
		return nil, nil, errors.New(ParseError)
	}

	// Split into seperate nodes (binary tree)
	n1, n2 := splitNode(n0, TokEqual)
	n3, n4 := splitNode(n1, TokComma)
	n5, n6 := splitNode(n2, TokPipe)
	n5, n7 := splitStackop(n5)

	// Construct form based on nodes
	sf.CurrentState = n3
	if sf.InputValue, err = strconv.Atoi(n4); err != nil {
		return nil, nil, err
	}
	of.Operation = n5
	if of.StackValue, err = strconv.Atoi(n7); err != nil {
		return nil, nil, err
	}
	of.NextState = n6
	return sf, of, nil
}

// Take some content, and transform into PdaMap
func Parse(rootStr string) (PdaMap, error) {
	ln := strings.Split(rootStr, TokNewLine)
	pm := make(PdaMap, len(ln))
	for _, v := range ln {
		if sf, of, err := parseLine(v); err != nil && err.Error() == ParseError {
			continue
		} else if err != nil && err.Error() != ParseError {
			return nil, err
		} else {
			pm[*sf] = *of
		}
	}
	return pm, nil
}

func main() {
	// Demo
	var rootStr = "A,4=^2|B\nB,2=!4|A\n"
	if pm, err := Parse(rootStr); err != nil {
		panic(err)
	} else {
		fmt.Printf("%v\n", pm)
	}
}
