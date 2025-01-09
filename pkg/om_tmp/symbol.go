package grammar

import "asritha.dev/compiler/pkg/utils"



type symbolType int

// Note: even if a Token is empty (epsilon), it must be passed to the parser.
// Note: end of file should also be passed as a Token matching the endOfFile var
const (
	NonTerm symbolType = iota
	Terminal
	Token
)

var (
	Epsilon   = symbol{sType: Terminal, data: ""}
	EndOfFile = symbol{sType: Token, data: "EOF"}
)

/*
Represents a single symbol (can be either a non-terminal or a terminal/token)
*/
type symbol struct {
	sType symbolType
	data  string
}

func NewSymbol(sType symbolType, data string) *symbol {
	return &symbol{
		sType: sType,
		data:  data,
	}
}

func (s symbol) String() string {
	return s.data
}

func (s symbol) Hash() int {
	return utils.HashStr(s.data) + int(s.sType)
}

func (s symbol) Equal(other symbol) bool {
	return s.sType == other.sType && s.data == other.data
}