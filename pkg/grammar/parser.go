package grammar

import (
	"fmt"
)

var (
	accept = action{accept: true}
)

type action struct {
	shift *lrAutomationState
	reduce *rule
	accept bool
}

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

type inputSymbol struct {
	symbol
	literal string
}

func (p parser) Parse(input []inputSymbol) error {
	stack := []*lrAutomationState{p.kernel}

	for {
		nextAction, ok := p.actionTable[stack[len(stack)-1]][input[0].symbol]
		if !ok {
			return fmt.Errorf("unexpected input")
		}

		if nextAction.accept {
			println("Complete!")
			return nil
		}

		if nextAction.shift != nil {
			fmt.Printf("shift %d\n", nextAction.shift.id)
			stack = append(stack, nextAction.shift)
			input = input[1:]
			continue
		}

		fmt.Printf("reduce %s\n", nextAction.reduce.String())

		//Pop states corresponding to β from the stack
		stack = stack[0:len(stack) - len(nextAction.reduce.sententialForm)]
		stack = append(stack, p.gotoTable[stack[len(stack)-1]][*NewNonTerm(nextAction.reduce.NonTerm)])
	}
}