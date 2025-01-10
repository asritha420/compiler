package parser

import "asritha.dev/compiler/pkg/grammar"

type LLTableParser struct {
	//stack -> TODO: make generic stack in utils folder
	parseTable map[string]string
	*grammar.Grammar
}

// parseTable is used to determine which rule should be applied for any combination of a non-terminal on the stack and next token in the input stream. An LL1 grammar will have only one rule to be applied for each combination

/*
LL(1) Parse Table Construction:
# TODo
 */

func (tp *LLTableParser) constructParseTable() {
	for _, rule := range tp.Rules {
		firstSet := tp.
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
