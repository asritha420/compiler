package scannergenerator

import (
	_ "regexp/syntax"
)

// TODO: figure out error handling in scanner
type TokenType int

type Token struct {
	TokenType
	Lexeme string
	// TODO: need literal?
	Line int
}
