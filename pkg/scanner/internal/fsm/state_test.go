package fsm

import (
	"reflect"
	"testing"

	"asritha.dev/compiler/pkg/utils"
)



func TestNewState(t *testing.T) {
	var id uint = 0
	s0 := NewNFAState(&id, false)
	s1 := NewNFAState(&id, true)
	s2 := NewNFAState(&id, false)

	utils.AssertEqual(t, "id", 3, id)

	utils.AssertEqual(t, "s0.id", 0, s0.id)
	utils.AssertEqual(t, "s0.IsAccepting", false, s0.IsAccepting)

	utils.AssertEqual(t, "s1.id", 1, s1.id)
	utils.AssertEqual(t, "s1.IsAccepting", true, s1.IsAccepting)

	utils.AssertEqual(t, "s2.id", 2, s2.id)
	utils.AssertEqual(t, "s2.IsAccepting", false, s2.IsAccepting)
}

func TestStateToEdge(t *testing.T) {
	var id uint = 0

	s0 := NewNFAState(&id, false)
	s1 := NewNFAState(&id, true)
	s2 := NewNFAState(&id, false)

	s0.AddTransition('a', s1, s2)
	s1.AddTransition('b', s0, s2)
	s2.AddTransition('c', s0, s1)

	edges := []*Edge{
		&Edge{s0, 'a', s1},
		&Edge{s0, 'a', s2},
		&Edge{s1, 'b', s0},
		&Edge{s1, 'b', s2},
		&Edge{s2, 'c', s0},
		&Edge{s2, 'c', s1},
	}
	t.Log(reflect.DeepEqual(edges, s0.getEdges()))
}

func TestSimple(t *testing.T) {
	var id uint = 0
	s1 := NewNFAState(&id, false)
	s2 := NewNFAState(&id, false)
	s3 := NewNFAState(&id, true)
	s4 := NewNFAState(&id, false)
	s5 := NewNFAState(&id, false)

	s1.AddTransition('a', s2, s4, s1)
	s2.AddTransition('b', s5, s3)
	s3.AddTransition('c', s1, s2)
	s4.AddTransition('d', s1)
	s5.AddTransition('e', s5)

	t.Log(s1)
}
