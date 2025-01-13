package parser

import (
	"fmt"
	"slices"

	. "asritha.dev/compiler/pkg/grammar"
	. "asritha.dev/compiler/pkg/scannergenerator"
	"asritha.dev/compiler/pkg/utils"
)

var (
	accept = action{accept: true}
)

type action struct {
	shift  *lrAutomationState
	reduce *Rule
	accept bool
}

func newReduce(r *Rule) action {
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
	actionTable map[*lrAutomationState]map[Symbol]action
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
		actionTable: make(map[*lrAutomationState]map[Symbol]action),
	}
	p.makeTables()

	return p
}

func (p parser) makeTables() {
	endAR := *NewAugmentedRule(p.FirstRule, len(p.FirstRule.SententialForm))

	for _, s := range p.states {
		p.gotoTable[s] = make(map[string]*lrAutomationState)
		p.actionTable[s] = make(map[Symbol]action)

		for ar, lookahead := range s.arLookaheadMap {
			nextSymbol := ar.getNextSymbol()
			if nextSymbol == nil {
				for symbol := range lookahead {
					p.actionTable[s][symbol] = newReduce(ar.Rule)
				}
				if ar == endAR {
					p.actionTable[s][EndOfInput] = accept
				}
				continue
			}

			if nextSymbol.SymbolType == TokenSymbol{
				p.actionTable[s][*nextSymbol] = newShift(s.transitions[*nextSymbol])
				continue
			}

			if nextSymbol.SymbolType == NonTermSymbol {
				p.gotoTable[s][nextSymbol.Name] = s.transitions[*nextSymbol]
			}

			// other symbol types should not appear but if they do don't do anything
		}
	}
}

type parseTreeNode struct {
	name     string
	literal  string
	children []*parseTreeNode
}

func newParseTreeNonTerm(name string, children []*parseTreeNode) *parseTreeNode {
	return &parseTreeNode{
		name:     name,
		children: children,
	}
}

func newParseTreeToken(t Token) *parseTreeNode {
	return &parseTreeNode{
		name:    t.Name,
		literal: t.Literal,
	}
}

func (p parser) Parse(input []Token) (*parseTreeNode, error) {
	stack := []*lrAutomationState{p.kernel}
	treeStack := make([]*parseTreeNode, 0)

	for {
		stackTop := stack[len(stack)-1] //top of stack
		var firstInput Token            // first input
		var firstInputSymbol Symbol     // first input as symbol (may be EndOfInput)
		if len(input) == 0 {
			firstInputSymbol = EndOfInput
		} else {
			firstInput = input[0]
			firstInputSymbol = *NewToken(firstInput.Name)
		}

		nextAction, ok := p.actionTable[stackTop][firstInputSymbol]
		if !ok {
			return nil, fmt.Errorf("unexpected input")
		}

		if nextAction.accept {
			// accept
			root := newParseTreeNonTerm(p.FirstRule.NonTerm, treeStack)
			return root, nil
		}

		if nextAction.shift != nil {
			// shift
			stack = append(stack, nextAction.shift)
			treeStack = append(treeStack, newParseTreeToken(firstInput))
			input = input[1:]
			continue
		}

		// reduce
		ruleLen := len(nextAction.reduce.SententialForm)
		newStackLen := len(stack) - ruleLen
		newTreeStackLen := len(treeStack) - ruleLen

		newNode := newParseTreeNonTerm(nextAction.reduce.NonTerm, slices.Clone(treeStack[newTreeStackLen:]))

		treeStack = treeStack[0:newTreeStackLen]
		treeStack = append(treeStack, newNode)

		stack = stack[0:newStackLen]
		stack = append(stack, p.gotoTable[stack[newStackLen-1]][nextAction.reduce.NonTerm])
	}
}

/*
Formats the parse tree by using the provided sets to remove extra nodes.

Note: order of ops goes removeSet > format children > compressSet > shorten
*/
func (node *parseTreeNode) Format(removeSet utils.Set[string], compressSet utils.Set[string], removeEmpty bool, shorten bool) *parseTreeNode {
	if _, ok := removeSet[node.name]; ok {
		return nil
	}

	for i := 0; i < len(node.children); i++ {
		child := node.children[i].Format(removeSet, compressSet, removeEmpty, shorten)
		if child == nil {
			node.children = utils.Remove(node.children, i)
			i--
		} else {
			node.children[i] = child
		}
	}

	if node.literal == "" && len(node.children) == 0 {
		return nil
	}

	if _, ok := compressSet[node.name]; ok {
		node.literal = node.GetLiteral()
		node.children = node.children[:0]
		return node
	}

	if shorten && len(node.children) == 1 {
		return node.children[0]
	}

	return node
}

/*
Recursively constructs the literal represented by a node
*/
func (node *parseTreeNode) GetLiteral() string {
	if len(node.children) == 0 {
		return node.literal
	}
	out := ""
	for _, n := range node.children {
		out += n.GetLiteral()
	}
	return out
}

func (p parser) MakeGraph(graphviz bool) string {
	if graphviz {
		return makeGraphviz(p.states)
	}
	return makeMermaid(p.states)
}
