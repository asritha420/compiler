package regex

import (
	"testing"

	"asritha.dev/compiler/pkg/scanner/internal/fsm"
)

func TestSimple(t *testing.T) {
	a := NewConst('a')
	b := NewConst('b')
	c := NewConst('c')

	aAltb := NewAlternation(a, b)
	aAltbKleen := NewKleeneStar(aAltb)
	caAltbKleen := NewConcatenation(c, aAltbKleen)

	var id uint = 0
	start, _, err := caAltbKleen.convertToNFA(&id)
	if err != nil {
		t.Fatal(err)
	}

	//equivalent NFA
	start2, end2, _ := caAltbKleen.convertToNFA(&id)
	end2.AddTransition(fsm.Epsilon, fsm.NewNFAState(&id, true))

	t.Log(start.IsEqual(start2))
}