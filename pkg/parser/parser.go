package parser

import (
	"fmt"
	"slices"

	"asritha.dev/compiler/pkg/grammar"
	"asritha.dev/compiler/pkg/scanner"
)

var (
	accept = action{accept: true}
)

type action struct {
	shift  *lrAutomationState
	reduce *grammar.Rule
	accept bool
}

func newReduce(r *grammar.Rule) action {
	return action{reduce: r}
}

func newShift(s *lrAutomationState) action {
	return action{shift: s}
}

type parser struct {
	*grammar.Grammar
	useLALR     bool
	kernel      *lrAutomationState
	states      []*lrAutomationState
	gotoTable   map[*lrAutomationState]map[string]*lrAutomationState
	actionTable map[*lrAutomationState]map[grammar.Symbol]action
}

func NewParser(g *grammar.Grammar, useLALR bool) *parser {
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
		actionTable: make(map[*lrAutomationState]map[grammar.Symbol]action),
	}
	p.makeTables()

	return p
}

func (p *parser) makeTables() {
	endAR := NewAugmentedRule(p.FirstRule, len(p.FirstRule.SententialForm))

	for _, s := range p.states {
		p.gotoTable[s] = make(map[string]*lrAutomationState)
		p.actionTable[s] = make(map[grammar.Symbol]action)

		for ar, lookahead := range s.arLookaheadMap {
			nextSymbol := ar.getNextSymbol()
			if nextSymbol == grammar.Epsilon {
				for symbol := range lookahead {
					p.actionTable[s][symbol] = newReduce(ar.Rule)
				}
				if ar == endAR {
					p.actionTable[s][grammar.EndOfInput] = accept
				}
				continue
			}

			if nextSymbol.SymbolType == grammar.TokenSymbol {
				p.actionTable[s][nextSymbol] = newShift(s.transitions[nextSymbol])
				continue
			}

			if nextSymbol.SymbolType == grammar.NonTermSymbol {
				p.gotoTable[s][nextSymbol.Name] = s.transitions[nextSymbol]
			}

			// other symbol types should not appear but if they do don't do anything
		}
	}
}

func (p parser) Parse(input []scanner.Token) (ParseTreeNode, error) {
	stack := []*lrAutomationState{p.kernel}
	treeStack := make([]ParseTreeNode, 0)

	for {
		stackTop := stack[len(stack)-1]     //top of stack
		var firstInput scanner.Token        // first input
		var firstInputSymbol grammar.Symbol // first input as symbol (may be EndOfInput)
		if len(input) == 0 {
			firstInputSymbol = grammar.EndOfInput
		} else {
			firstInput = input[0]
			firstInputSymbol = grammar.NewToken(firstInput.Name)
		}

		nextAction, ok := p.actionTable[stackTop][firstInputSymbol]
		if !ok {
			return nil, fmt.Errorf("unexpected input")
		}

		if nextAction.accept {
			// accept
			root := newParseTreeNonTerm(p.FirstRule, treeStack)
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
		ruleLen := nextAction.reduce.Len()
		newStackLen := len(stack) - ruleLen
		newTreeStackLen := len(treeStack) - ruleLen

		newNode := newParseTreeNonTerm(nextAction.reduce, slices.Clone(treeStack[newTreeStackLen:]))

		treeStack = treeStack[0:newTreeStackLen]
		treeStack = append(treeStack, newNode)

		stack = stack[0:newStackLen]
		stack = append(stack, p.gotoTable[stack[newStackLen-1]][nextAction.reduce.NonTerm])
	}
}

func (p parser) MakeGraph(graphviz bool) string {
	if graphviz {
		return makeGraphviz(p.states)
	}
	return makeMermaid(p.states)
}
