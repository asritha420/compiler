package fsm

import (
	"fmt"
	"strings"
)

func makeMermaidIdString(state *NFAState) string {
	id := fmt.Sprintf("id%d", state.id)
	if state.IsAccepting {
		id += fmt.Sprintf("(((%d)))", state.id)
	} else {
		id += fmt.Sprintf("((%d))", state.id)
	}
	return id
}

func MakeMermaid(rootState *NFAState) string {
	output := make([]string, 0)
	for _, edge := range rootState.getEdges() {
		output = append(output, edge.String())
	}
	output = append(output, fmt.Sprintf("START:::hidden -- start --> %s", makeMermaidIdString(rootState)))
	return strings.Join(output, "\n")
}
