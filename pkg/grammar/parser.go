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
		}
	}
}

// TODO may not need symbol type
type inputSymbol struct {
	symbol
	literal string
}

type parseTreeNode struct {
	name     string
	data     string
	children []*parseTreeNode
}

func newParseTreeNode(name string, data string, children []*parseTreeNode) *parseTreeNode {
	return &parseTreeNode{
		name:     name,
		data:     data,
		children: children,
	}
}

func (p parser) Parse(input []inputSymbol) (*parseTreeNode, error) {
	stack := []*lrAutomationState{p.kernel}
	treeStack := make([]*parseTreeNode, 0)

	for {
		s := stack[len(stack)-1] //top of stack
		a := input[0]            // first input

		nextAction, ok := p.actionTable[s][a.symbol]
		if !ok {
			return nil, fmt.Errorf("unexpected input")
		}

		if nextAction.accept {
			println("Complete!")
			root := newParseTreeNode(p.firstRule.NonTerm, p.firstRule.String(), treeStack)
			return root, nil
		}

		if nextAction.shift != nil {
			fmt.Printf("Shift %d\n", nextAction.shift.id)
			stack = append(stack, nextAction.shift)
			treeStack = append(treeStack, newParseTreeNode(a.name, a.literal, nil))
			input = input[1:]
			continue
		}

		fmt.Printf("Reduce %s\n", nextAction.reduce.String())

		//Pop states corresponding to β from the stack
		ruleLen := len(nextAction.reduce.sententialForm)
		newStackIdx := len(stack) - ruleLen
		newTreeStackIdx := len(treeStack) - ruleLen

		newNode := newParseTreeNode(nextAction.reduce.NonTerm, nextAction.reduce.String(), slices.Clone(treeStack[newTreeStackIdx:]))

		treeStack = treeStack[0:newTreeStackIdx]
		treeStack = append(treeStack, newNode)

		stack = stack[0:newStackIdx]
		stack = append(stack, p.gotoTable[stack[newStackIdx-1]][nextAction.reduce.NonTerm])
	}
}
