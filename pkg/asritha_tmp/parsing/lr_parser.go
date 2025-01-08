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

func NewShiftReduceParser(g *LRGrammar) *ShiftReduceParser {
	srp := &ShiftReduceParser{
		stack:     make(stack, 0),
		LRGrammar: g,
		indexer:   indexer{},
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
