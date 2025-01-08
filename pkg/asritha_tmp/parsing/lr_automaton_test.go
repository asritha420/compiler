package aparsing

import (
	"reflect"
	"testing"
)

// TODO: look up best practice to have test data in go
// grammar rules
var (
	pRule  = NewRule('P', "E")
	eRule  = NewRule('E', "E+T")
	eRule2 = NewRule('E', "T")
	tRule  = NewRule('T', "i(E)")
	tRule2 = NewRule('T', "i")
)

// LR(0) Automaton States
var (
	acceptState = &State{}
	state0      = &State{
		items: []*StateItem{
			{
				Rule:              pRule,
				dotIsToTheRightOf: 0,
			},
			{
				Rule:              eRule,
				dotIsToTheRightOf: 0,
			},
			{
				Rule:              eRule2,
				dotIsToTheRightOf: 0,
			},
			{
				Rule:              tRule,
				dotIsToTheRightOf: 0,
			},
		},
		transitions: map[rune]*State{
			'E': state1,
			'i': state4,
			'T': state8,
		},
	}

	state1 = &State{
		items: []*StateItem{
			{
				Rule:              pRule,
				dotIsToTheRightOf: 1,
			},
			{
				Rule:              eRule,
				dotIsToTheRightOf: 1,
			},
		},
		transitions: map[rune]*State{
			'$': acceptState, //$ be in valid terminals??
			'+': state2,
		},
	}

	state2 = &State{
		items: []*StateItem{
			// kernel
			{
				Rule:              eRule,
				dotIsToTheRightOf: 2,
			},
			// closure
			{
				Rule:              tRule,
				dotIsToTheRightOf: 0,
			},
			{
				Rule:              tRule2,
				dotIsToTheRightOf: 0,
			},
		},
		transitions: map[rune]*State{
			'T': state3,
			'i': state4,
		},
	}
	state3 = &State{}
	state4 = &State{}
	state5 = &State{}
	state6 = &State{}
	state7 = &State{}
	state8 = &State{}
)

var (
	grammar10 = LRGrammar{
		rules: []*Rule{
			pRule,
			eRule,
			eRule2,
			tRule,
			tRule2,
		},
		nonTerminals: []rune{'P', 'E', 'T'},
		terminals:    []rune{'+', 'i', '(', ')'},
	}
)

func TestNewState(t *testing.T) {
	shouldEqual := &State{
		items: []*StateItem{
			{
				Rule:              pRule,
				dotIsToTheRightOf: 0,
			},
			{
				Rule:              eRule,
				dotIsToTheRightOf: 0,
			},
			{
				Rule:              eRule2,
				dotIsToTheRightOf: 0,
			},
			{
				Rule:              tRule,
				dotIsToTheRightOf: 0,
			},
			{
				Rule:              tRule2,
				dotIsToTheRightOf: 0,
			},
		},
	}

	kernel := &StateItem{
		Rule:              pRule,
		dotIsToTheRightOf: 0,
	}
	if !reflect.DeepEqual(shouldEqual, grammar10.NewState(kernel)) {
		t.Errorf("FAILED")
	}
}
