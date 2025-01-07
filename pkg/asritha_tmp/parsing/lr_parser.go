package aparsing

type indexer struct {
	inputStream []rune
	curr        int // current UNCONSUMED index within inputStream
}

type ShiftReduceParser struct {
	stack
	*LRGrammar
	indexer
}

// LR0Automaton represents all possible rules currently under consideration by a shit-reduce parser
type LR0Automaton struct {
}

type StateRule struct {
	*Rule
	// dot indicates the parser's current position in the rule
	dotIsToTheRightOf int
}

type State struct {
	rules       []*StateRule
	transitions map[rune]*State
}

func NewState() *State {

}

func NewShiftReduceParser(g *LRGrammar) *ShiftReduceParser {
	srp := &ShiftReduceParser{
		stack:     make(stack, 0),
		LRGrammar: g,
		indexer:   indexer{},
	}

	state0 := &State{
		// creates kernel of state
		rules: []*StateRule{
			{
				Rule:              g.rules[0],
				dotIsToTheRightOf: 0,
			},
		},
	}

	return srp
}

// shift() consumes one token from input stream and pushes it onto the stack
func (srParser *ShiftReduceParser) shift() {
	srParser.stack.push(srParser.inputStream[srParser.curr])
	srParser.curr++
}

// reduce() uses a grammar rule "A -> a" to replace the sentential form "a" on stack with non-terminal "A"
func (srParser *ShiftReduceParser) reduce() {

}
