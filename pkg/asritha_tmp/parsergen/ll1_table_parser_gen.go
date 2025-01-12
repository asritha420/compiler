package aparsergen

// ParseTable determines which rule should be applied for any combination on non-terminal on the stack and next token on the input stream
type ParseTable struct {
}

type TableParserGen struct {
	*LLGrammar // an LL(1) grammar should have only one rule apply for each combination
	*ParseTable
	stack
}
