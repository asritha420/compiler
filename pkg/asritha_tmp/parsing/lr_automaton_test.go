package aparsing

import (
	"fmt"
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
			// kernel
			{
				Rule:              pRule,
				dotIsToTheRightOf: 0,
			},

			// non-kernel
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
			// kernel
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
			// non-kernel
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

	state3 = &State{
		items: []*StateItem{
			// kernel
			{
				Rule:              eRule,
				dotIsToTheRightOf: 3,
			},
		},
	}
	state4 = &State{
		items: []*StateItem{
			// kernel
			{
				Rule:              tRule,
				dotIsToTheRightOf: 1,
			},
			{
				Rule:              tRule2,
				dotIsToTheRightOf: 2,
			},
		},
		transitions: map[rune]*State{
			'(': state5,
		},
	}
	state5 = &State{
		items: []*StateItem{
			// kernel
			{
				Rule:              tRule,
				dotIsToTheRightOf: 3,
			},
			// non-kernel items
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
		transitions: map[rune]*State{
			'i': state4,
			'T': state8,
			'E': state6,
		},
	}
	state6 = &State{
		items: []*StateItem{
			// kernel
			{
				Rule:              tRule,
				dotIsToTheRightOf: 4,
			},
			{
				Rule:              eRule,
				dotIsToTheRightOf: 1,
			},
		},
		transitions: map[rune]*State{
			')': state7,
		},
	}
	state7 = &State{
		items: []*StateItem{
			// kernel
			{
				Rule:              tRule,
				dotIsToTheRightOf: 4,
			},
		},
	}
	state8 = &State{
		items: []*StateItem{
			// kernel
			{
				Rule:              eRule2,
				dotIsToTheRightOf: 1,
			},
		},
	}
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

	kernel := []*StateItem{
		{
			Rule:              pRule,
			dotIsToTheRightOf: 0,
		},
	}

	if !reflect.DeepEqual(shouldEqual, grammar10.NewState(kernel)) {
		t.Errorf("FAILED")
	}
}

func TestMermaidString(t *testing.T) {
	startState := &State{
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

	fmt.Println(startState.printInMermaid(0))
}
