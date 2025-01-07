package aparsing

// TODO: create newStateItem() ?? -> when to have Newxxx() and when to not?
type StateItem struct {
	*Rule
	// dot indicates the parser's current position in the rule
	dotIsToTheRightOf int
}

type State struct {
	items       []*StateItem
	transitions map[rune]*State
}

// NewLR0Automaton creates a LR0Automaton, which represents all possible rules currently under consideration by a shit-reduce parser
func (g *LRGrammar) NewLR0Automaton() *State {
	kernel := &StateItem{
		Rule:              g.rules[0],
		dotIsToTheRightOf: 0,
	}
	state0 := NewState(kernel)
	return state0
}

func NewState(kernel *StateItem) *State {
	s := &State{
		items: []*StateItem{kernel},
	}
	s.getClosure()
	return s
}

// TODO: shouldn't pass in grammar
func (s *State) getClosure(g *LRGrammar) {
	for _, item := range s.items {
		// if dot is to the right of a non-terminal
		symbol := item.Rule.production[item.dotIsToTheRightOf+1]
		if isValidNonTerminal(symbol) {
			_ = getClosureFor(g, symbol)
		}
	}
}

func getClosureFor(g *LRGrammar, nT rune) []*StateItem {
	startStateItems := make([]*StateItem, 0)
	// add all rules in the grammar that have symbol as a non-terminal
	for _, rule := range g.rules {
		if rule.nonTerminal == nT {
			startStateItems = append(startStateItems, &StateItem{
				Rule:              rule,
				dotIsToTheRightOf: -1,
			})
			startStateItems = append(startStateItems, getClosureFor(g, nT)...)
		}
	}
	return startStateItems
}

// should be a method on grammar or symbol type?
//func isValidNonTerminal(symbol rune) bool {
//
//}
