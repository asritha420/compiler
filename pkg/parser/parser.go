package parser

import "asritha.dev/compiler/pkg/grammar"

type LLTableParser struct {
	//stack -> TODO: make generic stack in utils folder
	parseTable []map[string]int // parseTable T[A, a], where A is a non-terminal, and a is a terminal; is used to determine which rule should be applied for any combination of a non-terminal on the stack and next token in the input stream
	// should be of fixed length
	*grammar.Grammar
}

/*
LL(1) Parse Table Construction:

*/
func (tp *LLTableParser) constructParseTable() {
	for i, rule := range tp.Rules {
		firstSet := tp.FirstSets[rule.NonTerm]
		for k, _ := range firstSet {
			tp.parseTable[i][k.String()] = i
		}
	}
}

func (tp *LLTableParser) Parse() {

}

func NewLLTableParser(g *grammar.Grammar) *LLTableParser {

}

/*
user will:

tokens := []Token // from scanner
grammar := NewGrammar(xxxx)
tableParser := NewLLTableParser(grammar)
parseTree := tableParser.parse(tokens)
*/
