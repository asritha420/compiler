package parsergen

type Rule struct {
	NonTerminal string
	Productions []string
	firstSet    []string
	followSet   []string
}

func NewRule(nT string, productions []string) *Rule {
	return &Rule{
		NonTerminal: nT,
		Productions: productions,
		firstSet: make([]string, 0),
		followSet: make([]string, 0),
	}
}
