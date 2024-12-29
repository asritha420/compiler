package scanner

type TokenType int

type Token struct {
	TType   TokenType
	Lexeme  string
	LineNum int
}

type TokenInitInfo struct {
	TType       TokenType
	Name, Regex string
	//fa *State
}
