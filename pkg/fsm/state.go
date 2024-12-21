package fsm

const (
	Epsilon = 0x00
)

type StateEdge struct {
	transition rune
	next       State
}

type State interface {
	GetId() uint
	GetEdges() []StateEdge
	IsAccepting() bool
}

type NFAState struct {
	id          uint
	transitions map[rune][]*NFAState
	isAccepting bool
}

func NewNFAState(id *uint) *NFAState {
	state := &NFAState{
		id: *id,
		transitions: make(map[rune][]*NFAState),
	}
	*id++
	return state
}

func (state *NFAState) GetId() uint {
	return state.id
}

func (state *NFAState) AddTransition(transition rune, newStates ...*NFAState) {
	if _, ok := state.transitions[transition]; !ok {
		state.transitions[transition] = make([]*NFAState, len(newStates))
	} 

	state.transitions[transition] = append(state.transitions[transition], newStates...)
}

func (state *NFAState) SetAccepting(accepting bool) {
	state.isAccepting = accepting
}

func (state *NFAState) GetEdges() []StateEdge {
	out := make([]StateEdge, 0)
	for transition, nextStates := range state.transitions {
		for _, nextState := range nextStates {
			out = append(out, StateEdge{
				transition: transition,
				next:       nextState,
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
			transition: transitions,
			next:       nextState,
		})

	}
	return out
}

func (state *DFAState) IsAccepting() bool {
	return state.isAccepting
}
