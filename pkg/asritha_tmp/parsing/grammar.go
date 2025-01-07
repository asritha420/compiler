package aparsing

type Rule struct {
	nonTerminal rune
	production  []rune
}

type LRGrammar struct {
	rules                   []*Rule
	nonTerminals, terminals []rune
}
