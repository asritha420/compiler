package grammar

type symbolType int

const (
	nonTerminal symbolType = iota
	terminal
	terminalLowercaseRange
	terminalUppercaseRange
	terminalNumberRange
	epsilon
)

func (s symbolType) String() string {
	return [...]string{"nonTerminal", "terminal", "terminalLowercaseRange", "terminalUppercaseRange", "terminalNumberRange", "epsilon"}[s]
}

type symbol struct {
	symbolType    symbolType
	validLiterals []string //replace with isValid function? 
}

func newNonTerminalSymbol(literal string) *symbol {
	return &symbol{
		symbolType:    nonTerminal,
		validLiterals: []string{literal},
	}
}

func newTerminalSymbol(literal string) *symbol {
	symbolType := terminal

	if literal == " " {
		symbolType = epsilon
	}

	return &symbol{
		symbolType:    symbolType,
		validLiterals: []string{literal},
	}
}

func newTerminalRangeSymbol(rangeType symbolType) *symbol {
	var validLiterals []string
	var from rune
	var to rune

	switch rangeType {
	case terminalLowercaseRange:
		from = 'a'
		to = 'z'
	case terminalUppercaseRange:
		from = 'A'
		to = 'Z'
	case terminalNumberRange:
		from = 0
		to = 9
	}

	for l := from; l <= to; l++ {
		validLiterals = append(validLiterals, string(l))
	}

	return &symbol{
		symbolType:    rangeType,
		validLiterals: validLiterals,
	}
}
