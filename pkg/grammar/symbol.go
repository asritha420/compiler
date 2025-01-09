package grammar

import "asritha.dev/compiler/pkg/utils"

type symbolType int
type TokenType int

// Note: even if a Token is empty (epsilon), it must be passed to the parser.
// Note: end of file should also be passed as a Token matching the endOfFile var
// 
const (
	NonTerm symbolType = iota
	Token
	Epsilon
	EndOfInput
)

var (
	EpsilonSymbol   = symbol{symbolType: Epsilon}
	EndOfInputSymbol = symbol{symbolType: EndOfInput}
)

/*
Represents a single symbol (can be either a non-terminal or a terminal/token)
*/
type symbol struct {
	symbolType
	TokenType
	data string
}

func (s symbol) Hash() int {
	switch s.symbolType {
	case NonTerm:
		return utils.HashStr(s.data)
	case Token:
		return int(s.TokenType) + 2
	case Epsilon:
		return 0
	case EndOfInput:
		return 1
	}
	return -1
}

func (s symbol) Equal(other symbol) bool {
	if s.symbolType != other.symbolType {
		return false
	}

	if s.symbolType == EndOfInput || s.symbolType == Epsilon {
		return true
	}

	if s.symbolType == NonTerm {
		return s.data == other.data
	}

	return s.TokenType == other.TokenType
}

func NewNonTerm(name string) *symbol {
	return &symbol{
		symbolType: NonTerm,
		data: name,
	}
}

func NewToken(tokenType TokenType, data string) *symbol {
	return &symbol{
		symbolType: Token,
		TokenType: tokenType,
		data: data,
	}
}

//TODO
func (s symbol) String() string {
	return ""
}