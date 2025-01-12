package fsm

import (
	"fmt"
	"strings"
)

const (
	Epsilon = 0x00
)

type NFAState struct {
	id          uint
	transitions map[rune][]*NFAState
	IsAccepting bool
}

type Edge struct {
	start      *NFAState
	transition rune
	end        *NFAState
}

func (state *NFAState) getEdgesRecursion(edges []*Edge, closed map[uint]struct{}) []*Edge {
	id := state.id
	if _, ok := closed[id]; ok {
		return edges
	}
	closed[id] = struct{}{}
	for transition, nextStates := range state.transitions {
		for _, nextState := range nextStates {
			edges = append(edges, &Edge{start: state, transition: transition, end: nextState})
			edges = nextState.getEdgesRecursion(edges, closed)
		}
	}
	return edges
}

func (state *NFAState) getEdges() []*Edge {
	return state.getEdgesRecursion(make([]*Edge, 0), make(map[uint]struct{}))
}

func makeMermaidIdString(state *NFAState) string {
	id := fmt.Sprintf("id%d", state.id)
	if state.IsAccepting {
		id += fmt.Sprintf("(((%d)))", state.id)
	} else {
		id += fmt.Sprintf("((%d))", state.id)
	}
	return id
}

func (e *Edge) String() string {
	return fmt.Sprintf("%s -- %c --> %s", makeMermaidIdString(e.start), e.transition, makeMermaidIdString(e.end))
}

func (state *NFAState) String() string {
	output := make([]string, 0)
	for _, edge := range state.getEdges() {
		output = append(output, edge.String())
	}
	output = append(output, fmt.Sprintf("START:::hidden -- start --> %s", makeMermaidIdString(state)))
	return strings.Join(output, "\n")
}

func NewNFAState(id *uint, isAccepting bool) *NFAState {
	state := &NFAState{
		id:          *id,
		transitions: make(map[rune][]*NFAState),
		IsAccepting: isAccepting,
	}
	*id++
	return state
}

func (state *NFAState) AddTransition(transition rune, newStates ...*NFAState) {
	if _, ok := state.transitions[transition]; !ok {
		state.transitions[transition] = make([]*NFAState, 0)
	}

	state.transitions[transition] = append(state.transitions[transition], newStates...)
}

func (state *NFAState) isStateEqual(other *NFAState, idMap map[uint]uint) bool {
	if otherId, ok := idMap[state.id]; ok {
		if other.id != otherId {
			return false
		}
	} else {
		if state.IsAccepting != other.IsAccepting {
			return false
		}
		idMap[state.id] = other.id
	}
	return true
}

//TODO make order not matter
/*
Check if 2 NFAs are the same. This does NOT check id but rather if the 2 states will produce the same graph.
*/
func (state *NFAState) IsEqual(other *NFAState) bool {
	// first check if they are the same
	if state == other {
		return true
	}

	// get all edges
	edges := state.getEdges()
	otherEdges := other.getEdges()

	// check for same number of edges
	if len(edges) != len(otherEdges) {
		return false
	}

	idMap := make(map[uint]uint)
	for i, edge := range edges {
		otherEdge := otherEdges[i]
		//check idMap
		if !edge.start.isStateEqual(otherEdge.start, idMap) {
			return false
		}

		if !edge.end.isStateEqual(otherEdge.end, idMap) {
			return false
		}
	}

	return true
}
