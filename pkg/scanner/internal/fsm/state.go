package fsm

import "fmt"

const (
	Epsilon = 0x00
)

type NFAState struct {
	id          uint
	transitions map[rune][]*NFAState
	IsAccepting bool
}

type Edge struct {
	start *NFAState
	transition rune
	end *NFAState
}

func (e *Edge) String() string {
	return fmt.Sprintf("%s -- %c --> %s", makeMermaidIdString(e.start), e.transition, makeMermaidIdString(e.end))
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


func (s *NFAState) getEdgesRecursion(edges []*Edge, closed map[uint]struct{}) []*Edge {
	id := s.id
	if _, ok := closed[id]; ok {
		return edges
	}
	closed[id] = struct{}{}
	for transition, nextStates := range s.transitions {
		for _, nextState := range nextStates {
			edges = append(edges, &Edge{start: s, transition: transition, end: nextState})
			edges = nextState.getEdgesRecursion(edges, closed)
		}
	}
	return edges
}

func (s *NFAState) getEdges() []*Edge {
	return s.getEdgesRecursion(make([]*Edge, 0), make(map[uint]struct{}))
}

//TODO make order not matter
/*
Check if 2 NFAs are the same. This does NOT check id but rather if the 2 states will produce the same graph.
*/
func (s *NFAState) IsEqual(other *NFAState) bool {
	// first check if they are the same
	if s == other {
		return true
	}

	// get all edges
	edges := s.getEdges()
	otherEdges := other.getEdges()

	if len(edges) != len(otherEdges) {
		return false
	}

	idMap := make(map[uint]uint)
	for i, edge := range edges {
		otherEdge := otherEdges[i]
		//check idMap
		if otherId, ok := idMap[edge.start.id]; ok {
			if otherEdge.start.id != otherId {
				return false
			}
		} else {
			if edge.start.IsAccepting != otherEdge.start.IsAccepting {
				return false
			}
			idMap[edge.start.id] = otherEdge.start.id
		}
	}
}