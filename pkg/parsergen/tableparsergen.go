package parsergen

import (
	"asritha.dev/compiler/pkg/grammar"
)

// needs a grammar, parse table, stack
type TableParserGen struct {
	*grammar.Grammar                   // TODO: need to verify that this is LL1 grammar
	parseTable       map[string]string // should be symbol
	stack                              // create stack in utils
}
