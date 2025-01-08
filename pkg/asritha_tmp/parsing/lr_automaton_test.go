package aparsing

import (
	"reflect"
	"testing"
)

var (
	pRule     = NewRule('P', "E")
	eRule     = NewRule('E', "E+T")
	eRule2    = NewRule('E', "T")
	tRule     = NewRule('T', "i(E)")
	tRule2    = NewRule('T', "i")
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
