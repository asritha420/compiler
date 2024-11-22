package parsergen

type Rule struct {
	NonTerminal string
	Productions []string
	firstSet    []string
	followSet   []string
}

func NewRule(nT string, productions []string) (*Rule, error) {
	//TODO: implement error handling if inaccurate rule
	r := &Rule{
		NonTerminal: nT,
		Productions: productions,
	}
	return r, nil
}
