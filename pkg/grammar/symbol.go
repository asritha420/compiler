package grammar

import (
	"reflect"

	"asritha.dev/compiler/pkg/utils"
)

type symbolType int

// Note: even if a token is empty (epsilon), it must be passed to the parser.
// Note: end of file should also be passed as a token matching the endOfFile var
const (
	nonTerm symbolType = iota
	token
	epsilon
	endOfInput
)

var (
	Epsilon    = symbol{symbolType: epsilon}
	EndOfInput = symbol{symbolType: endOfInput}
)

/*
Represents a single symbol (can be either a non-terminal or a terminal/token)
*/
type symbol struct {
	symbolType
	name string
}

func (s symbol) Hash() int {
	return int(s.symbolType) + utils.HashStr(s.name)
}

func (s symbol) Equal(other symbol) bool {
	return reflect.DeepEqual(s, other)
}

func NewNonTerm(name string) *symbol {
	return &symbol{
		symbolType: nonTerm,
		name:       name,
	}
}

func NewToken(name string) *symbol {
	return &symbol{
		symbolType: token,
		name:       name,
	}
}

func (s symbol) String() string {
	switch s.symbolType {
	case epsilon:
		return "ε"
	case endOfInput:
		return "$"
	}
	return s.name
}
