package fsm

const (
	Epsilon = 0x00
)

type NFAState struct {
	Id          uint
	Transitions map[rune][]*NFAState
	IsAccepting bool
}

func NewNFAState(id *uint, isAccepting bool) *NFAState {
	state := &NFAState{
		Id:          *id,
		Transitions: make(map[rune][]*NFAState),
		IsAccepting: isAccepting,
	}
	*id++
	return state
}

func (state *NFAState) AddTransition(transition rune, newStates ...*NFAState) {
	if _, ok := state.Transitions[transition]; !ok {
		state.Transitions[transition] = make([]*NFAState, 0)
	}

	state.Transitions[transition] = append(state.Transitions[transition], newStates...)
}
