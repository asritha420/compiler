package grammar

import (
	"fmt"
	"slices"
)

var (
	accept = action{accept: true}
)

type action struct {
	shift  *lrAutomationState
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
	useLALR     bool
	kernel      *lrAutomationState
	states      []*lrAutomationState
	gotoTable   map[*lrAutomationState]map[string]*lrAutomationState
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
		Grammar:     g,
		useLALR:     useLALR,
		kernel:      kernel,
		states:      states,
		gotoTable:   make(map[*lrAutomationState]map[string]*lrAutomationState),
		actionTable: make(map[*lrAutomationState]map[symbol]action),
	}
	p.makeTables()

	return p
}

func (p parser) makeTables() {
	endAR := *NewAugmentedRule(p.Grammar.Rules[0], len(p.Grammar.Rules[0].sententialForm))

	for _, s := range p.states {
		p.gotoTable[s] = make(map[string]*lrAutomationState)
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
				p.gotoTable[s][nextSymbol.name] = s.transitions[*nextSymbol]
			}

			// other symbol types should not appear but if they do don't do anything
		}
	}
}

type parseTreeNode interface{}

type parseTreeNonTerm struct {
	name     string
	children []parseTreeNode
}

func newParseTreeNonTerm(name string, children []parseTreeNode) *parseTreeNonTerm {
	return &parseTreeNonTerm{
		name:     name,
		children: children,
	}
}

type Token struct {
	name    string
	literal string
}

func (p parser) Parse(input []Token) (parseTreeNode, error) {
	stack := []*lrAutomationState{p.kernel}
	treeStack := make([]parseTreeNode, 0)

	for {
		stackTop := stack[len(stack)-1] //top of stack
		var firstInput Token // first input
		var firstInputSymbol symbol             // first input as symbol (may be EndOfInput)
		if len(input) == 0 {
			firstInputSymbol = EndOfInput
		} else {
			firstInput := input[0]
			firstInputSymbol = *NewToken(firstInput.name)
		}

		nextAction, ok := p.actionTable[stackTop][firstInputSymbol]
		if !ok {
			return nil, fmt.Errorf("unexpected input")
		}

		if nextAction.accept {
			// accept
			root := newParseTreeNonTerm(p.firstRule.NonTerm, treeStack)
			return root, nil
		}

		if nextAction.shift != nil {
			// shift
			stack = append(stack, nextAction.shift)
			treeStack = append(treeStack, firstInput)
			input = input[1:]
			continue
		}

		// reduce
		ruleLen := len(nextAction.reduce.sententialForm)
		newStackLen := len(stack) - ruleLen
		newTreeStackLen := len(treeStack) - ruleLen

		newNode := newParseTreeNonTerm(nextAction.reduce.NonTerm, slices.Clone(treeStack[newTreeStackLen:]))

		treeStack = treeStack[0:newTreeStackLen]
		treeStack = append(treeStack, newNode)

		stack = stack[0:newStackLen]
		stack = append(stack, p.gotoTable[stack[newStackLen-1]][nextAction.reduce.NonTerm])
	}
}
