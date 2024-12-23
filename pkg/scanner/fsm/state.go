package scanner

type FAState struct {
	Id          uint
	IsAccepting bool
}

type NFAState struct {
	FAState     //TODO: make this a pointer?
	Transitions map[rune][]*NFAState
}

type DFAState struct {
	FAState
	transitions []*DFATransition // If a character in the alphabet does not have an explicit transition, it is assumed to lead to a dead state by default
}

type DFATransition struct {
	characterClass []rune
	pointingTo     *DFAState
}
