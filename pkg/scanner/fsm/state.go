package scanner

const (
	Epsilon = 0x00
)

type NFAState struct {
	id          uint
	transitions map[rune][]*NFAState
	accepting   bool
}

func NewNFAState(id *uint, accepting bool) *NFAState {
	state := &NFAState{
		id:          *id,
		transitions: make(map[rune][]*NFAState),
		accepting:   accepting,
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
