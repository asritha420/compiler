package fsm

const (
	Epsilon = 0x00
)

type NFAState struct {
	id          uint
	transitions map[rune][]*NFAState
	isAccepting bool
	isPseudoDFA bool
}

func NewNFAState(id *uint) *NFAState {
	state := &NFAState{
		id: *id,
		transitions: make(map[rune][]*NFAState),
	}
	*id++
	return state
}

func newPseudoDFAState(id *uint, accepting bool) *NFAState {
	state := &NFAState{
		id: *id,
		transitions: make(map[rune][]*NFAState),
		isAccepting: accepting,
		isPseudoDFA: true,
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

func (state *NFAState) IsAccepting() bool {
	return state.isAccepting
}

func (state *NFAState) IsPseudoDFA() bool {
	return state.isPseudoDFA
}