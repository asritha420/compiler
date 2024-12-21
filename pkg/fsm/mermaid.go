package fsm

import (
	"fmt"
	"strings"
)

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
		if edge.transition == Epsilon {
			transition = "ɛ"
		} else {
			transition = fmt.Sprintf("%c", edge.transition)
		}
		edges = append(edges, fmt.Sprintf("%s -- %s --> %s", makeMermaidIdString(rootState), transition, makeMermaidIdString(edge.next)))
		edges, closed = makeMermaidRecursion(edge.next, edges, closed)
	}
	return edges, closed
}

func MakeMermaid(rootState State) string {
	edges, _ := makeMermaidRecursion(rootState, make([]string, 0), make(map[uint]struct{}))
	edges = append(edges, fmt.Sprintf("START:::hidden -- start --> %s", makeMermaidIdString(rootState)))
	return strings.Join(edges, "\n")
}
