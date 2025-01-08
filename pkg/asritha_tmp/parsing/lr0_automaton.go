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
	state0 := g.NewState()
	return state0
}

func (g *LRGrammar) NewState() *State {
	s := &State{
		items: make([]*StateItem, 0),
	}
	//TODO: have comment saying g.nonTerminals[0] is kernel? -> also this function should be reusable ??
	s.getClosureFor(g, g.nonTerminals[0])
	return s
}

// TODO: have separate fields for kernelItems and nonKernelItems in state?
// TODO: don't like how I am passing in grammar
//
//	func (s *State) getClosure(g *LRGrammar) {
//		s.items = s.getClosureFor(g, g.nonTerminals[0])
//	} //
//
// start with P
func (s *State) getClosureFor(g *LRGrammar, nT rune) {
	// TODO: rename this nonKernelItems?
	//startStateItems := make([]*StateItem, 0)
	// add all rules in the grammar that have symbol as a non-terminal
	matchingRules := g.getRulesForNT(nT)
	for _, rule := range matchingRules {
		item := &StateItem{
			Rule:              rule,
			dotIsToTheRightOf: 0,
		}

		s.items = append(s.items, item)

		// if dot is to the right of a non-terminal
		symbol := item.Rule.production[item.dotIsToTheRightOf]
		if g.isValidNonTerminal(symbol) && symbol != nT {
			//startStateItems = append(startStateItems, s.getClosureFor(g, symbol)...)
			s.getClosureFor(g, symbol)
		}
	}
}
