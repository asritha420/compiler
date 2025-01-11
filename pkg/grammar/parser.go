package grammar

type inputSymbol struct {
	symbol
	literal string
}

type action struct {
	shift *lrAutomationState
	reduce *rule
	accept bool
}

var (
	accept = action{accept: true}
)

func newReduce(r *rule) action {
	return action{reduce: r}
}

func newShift(s *lrAutomationState) action {
	return action{shift: s}
}

type parser struct {
	*Grammar
	useLALR bool
	kernel *lrAutomationState
	states []*lrAutomationState
	gotoTable map[*lrAutomationState]map[symbol]*lrAutomationState
	actionTable map[*lrAutomationState]map[symbol]action
}

func NewParser(g *Grammar, useLALR bool) *parser {
	var kernel *lrAutomationState
	var states []*lrAutomationState
	if useLALR {
		kernel, states = generateLALR(g)
	} else {
		kernel, states = generateLR1(g)
	}
	print(makeMermaid(states))
	p := &parser{
		Grammar: g,
		useLALR: useLALR,
		kernel: kernel,
		states: states,
		gotoTable: make(map[*lrAutomationState]map[symbol]*lrAutomationState),
		actionTable: make(map[*lrAutomationState]map[symbol]action),
	}
	p.makeTables()

	return p
}

func (p parser) makeTables() {
	endAR := *NewAugmentedRule(p.Grammar.Rules[0], len(p.Grammar.Rules[0].sententialForm))

	for _, s := range p.states {
		p.gotoTable[s] = make(map[symbol]*lrAutomationState)
		p.actionTable[s] = make(map[symbol]action)

		for ar, lookahead := range s.arLookaheadMap {
			nextSymbol := ar.getNextSymbol()
			if nextSymbol == nil {
				for symbol := range lookahead {
					p.actionTable[s][symbol] = newReduce(ar.rule)
				}
				if ar == endAR {
					p.actionTable[s][EndOfInput] = accept
				}
				continue
			}

			if nextSymbol.symbolType == token {
				p.actionTable[s][*nextSymbol] = newShift(s.transitions[*nextSymbol])
				continue
			}

			if nextSymbol.symbolType == nonTerm {
				p.gotoTable[s][*nextSymbol] = s.transitions[*nextSymbol]
			}
		}
	}
}

// func (p parser) Parse(input []inputSymbol) {
// 	stack := []*lrAutomationState{p.kernel}
// 	for {

// 	}
// }