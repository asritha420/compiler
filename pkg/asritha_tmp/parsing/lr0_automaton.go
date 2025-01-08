package aparsing

// TODO: create newStateItem() ?? -> when to have Newxxx() and when to not?
type StateItem struct {
	*Rule
	// dot indicates the parser's current position in the rule
	dotIsToTheRightOf int //make required somehow, but it wont matter bc a user won't initalize this?
}

// TODO: have field tags for kernel vs non-kernel items
type State struct {
	items       []*StateItem
	transitions map[rune]*State
}

// TODO: figure out when to "reach accepting"
// NewLR0Automaton creates a LR0Automaton, which represents all possible rules currently under consideration by a shit-reduce parser
func (g *LRGrammar) NewLR0Automaton() *State {
	state0Kernel := []*StateItem{
		{
			Rule:              g.rules[0],
			dotIsToTheRightOf: 0,
		},
	}
	state0 := g.NewState(state0Kernel)
	return state0
}

// TODO: create mermaid graph printer

// TODO: if having same state, should point to same state instead of redoing it
func (g *LRGrammar) NewState(kernel []*StateItem) *State {
	s := &State{
		items: make([]*StateItem, 0),
	}
	for _, kernelItem := range kernel {
		s.getClosure(g, kernelItem)
	}

	s.CreateTransitions(g)
	return s
}

// create transititions for each symbol to the right of dot
func (s *State) CreateTransitions(g *LRGrammar) {
	// get list of all possible transitions
	transitionKernels := make(map[rune][]*StateItem)
	for _, item := range s.items {
		if item.dotIsToTheRightOf == len(item.production) {
			continue
		}
		transitionSymbol := item.production[item.dotIsToTheRightOf]

		newItem := &StateItem{
			Rule:              item.Rule,
			dotIsToTheRightOf: item.dotIsToTheRightOf,
		}

		if existing, ok := transitionKernels[transitionSymbol]; ok {
			existing = append(existing, newItem)
		} else {
			transitionKernels[transitionSymbol] = append(existing, newItem)
		}
	}

	for transition, kernel := range transitionKernels {
		s.transitions[transition] = g.NewState(kernel)
	}
}

// TODO: have separate fields for kernelItems and nonKernelItems in state?
func (s *State) getClosure(g *LRGrammar, stateItem *StateItem) {
	// add all rules in the grammar that have symbol as a non-terminal
	symbol := stateItem.production[stateItem.dotIsToTheRightOf]
	if g.isValidNonTerminal(symbol) && symbol != stateItem.nonTerminal {
		matchingRules := g.getRulesForNT(symbol) //TODO: have consistent naming everywhere for nT/nonTerminal
		for _, rule := range matchingRules {
			item := &StateItem{
				Rule:              rule,
				dotIsToTheRightOf: 0,
			}

			s.items = append(s.items, item)

			s.getClosure(g, item)
		}
	}
}
