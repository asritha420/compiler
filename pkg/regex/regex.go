package regex

type Regex string

func (r Regex) GetAST() (Node, error) {

	parser := regexParser{
		// TODO: fix this dep
		regex: []rune(r),
		curr:  0,
	}

	AST, err := parser.parse()

	if err != nil {
		return nil, err
	}

	return AST, nil
}
