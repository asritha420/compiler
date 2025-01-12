package grammar

import (
	"reflect"

	"asritha.dev/compiler/pkg/utils"
)

type symbolType int

// Note: even if a token is empty (epsilonSymbol), it must be passed to the parser.
// Note: end of file should also be passed as a token matching the endOfFile var
const (
	NonTermSymbol symbolType = iota
	TokenSymbol
	epsilonSymbol
	endOfInputSymbol
)

var (
	Epsilon    = Symbol{SymbolType: epsilonSymbol}
	EndOfInput = Symbol{SymbolType: endOfInputSymbol}
)

/*
Represents a single Symbol (can be either a non-terminal or a terminal/token)
*/
type Symbol struct {
	SymbolType symbolType
	Name       string
}

func (s Symbol) Hash() int {
	return int(s.SymbolType) + utils.HashStr(s.Name)
}

func (s Symbol) Equal(other Symbol) bool {
	return reflect.DeepEqual(s, other)
}

func NewNonTerm(name string) *Symbol {
	return &Symbol{
		SymbolType: NonTermSymbol,
		Name:       name,
	}
}

func NewToken(name string) *Symbol {
	return &Symbol{
		SymbolType: TokenSymbol,
		Name:       name,
	}
}

func (s Symbol) String() string {
	switch s.SymbolType {
	case epsilonSymbol:
		return "ε"
	case endOfInputSymbol:
		return "$"
	case TokenSymbol, NonTermSymbol:
		return s.Name
	}
	return ""
}
