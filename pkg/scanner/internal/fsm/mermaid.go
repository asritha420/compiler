package fsm

import (
	"fmt"
	"strings"
)

func makeMermaidIdString(state *NFAState) string {
	id := fmt.Sprintf("id%d", state.id)
	if state.accepting {
		id += fmt.Sprintf("(((%d)))", state.id)
	} else {
		id += fmt.Sprintf("((%d))", state.id)
	}
	return id
}

func makeMermaidRecursion(rootState *NFAState, edges []string, closed map[uint]struct{}) []string {
	id := rootState.id
	if _, ok := closed[id]; ok {
		return edges
	}
	closed[id] = struct{}{}
	for transition, nextStates := range rootState.transitions {
		if transition == Epsilon {
			transition = 'ɛ'
		}
		for _, nextState := range nextStates {
			edges = append(edges, fmt.Sprintf("%s -- %c --> %s", makeMermaidIdString(rootState), transition, makeMermaidIdString(nextState)))
			edges = makeMermaidRecursion(nextState, edges, closed)
		}
	}
	return edges
}

func MakeMermaid(rootState *NFAState) string {
	edges := makeMermaidRecursion(rootState, make([]string, 0), make(map[uint]struct{}))
	edges = append(edges, fmt.Sprintf("START:::hidden -- start --> %s", makeMermaidIdString(rootState)))
	return strings.Join(edges, "\n")
}
