package scannergen

import (
	"fmt"
	"strings"
)

type StateEdge struct {
	Transition rune
	Next       State
}

type State interface {
	GetId() uint
	GetEdges() []StateEdge
	IsAccepting() bool
}

type NFAState struct {
	id          uint
	transitions map[rune][]*NFAState //input:nfaState
	isAccepting bool
}

func (state *NFAState) GetId() uint {
	return state.id
}

func (state *NFAState) GetEdges() []StateEdge {
	out := make([]StateEdge, 0)
	for transition, nextStates := range state.transitions {
		for _, nextState := range nextStates {
			out = append(out, StateEdge{
				Transition: transition,
				Next:       nextState,
			})
		}
	}
	return out
}

func (state *NFAState) IsAccepting() bool {
	return state.isAccepting
}

type DFAState struct {
	id          uint
	transitions map[rune]*DFAState
	isAccepting bool
}

func (state *DFAState) GetId() uint {
	return state.id
}

func (state *DFAState) GetEdges() []StateEdge {
	out := make([]StateEdge, 0)
	for transitions, nextState := range state.transitions {
		out = append(out, StateEdge{
			Transition: transitions,
			Next:       nextState,
		})

	}
	return out
}

func (state *DFAState) IsAccepting() bool {
	return state.isAccepting
}

func makeMermaidIdString(state State) string {
	id := fmt.Sprintf("id%d", state.GetId())
	if state.IsAccepting() {
		id += fmt.Sprintf("(((%d)))", state.GetId())
	} else {
		id += fmt.Sprintf("((%d))", state.GetId())
	}
	return id
}

func makeMermaidRecursion(rootState State, edges []string, closed map[uint]struct{}) ([]string, map[uint]struct{}) {
	id := rootState.GetId()
	if _, ok := closed[id]; ok {
		return edges, closed
	}
	closed[id] = struct{}{}
	for _, edge := range rootState.GetEdges() {
		transition := ""
		if edge.Transition == epsilon {
			transition = "ɛ"
		} else {
			transition = fmt.Sprintf("%c", edge.Transition)
		}
		edges = append(edges, fmt.Sprintf("%s -- %s --> %s", makeMermaidIdString(rootState), transition, makeMermaidIdString(edge.Next)))
		edges, closed = makeMermaidRecursion(edge.Next, edges, closed)
	}
	return edges, closed
}

func MakeMermaid(rootState State) string {
	edges, _ := makeMermaidRecursion(rootState, make([]string, 0), make(map[uint]struct{}))
	edges = append(edges, fmt.Sprintf("START:::hidden -- start --> %s", makeMermaidIdString(rootState)))
	return strings.Join(edges, "\n")
}