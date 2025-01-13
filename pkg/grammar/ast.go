package grammar

type grammarAST interface {}

type kleenStar struct {
	child grammarAST
}

type optional struct {
	child grammarAST
}

type oneOrMore struct {
	
}