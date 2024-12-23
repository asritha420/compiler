package parser

type Grammar struct {
	Rules        []*Rule
	nonTerminals [][]rune
	terminals    []rune
}
