package scannergenerator

type TokenType int

type Token struct {
	TokenType
	Lexeme string
	Line   int
}
