package scanner

import (
	"asritha.dev/compiler/pkg/regex"
)

type Scanner interface {
	Scan()
}

// need to be list, not map
func New(map[string]regex.Regex) Scanner {

	return nil
}

type scannerGen struct {
}

//must implement parseTreeNode interface!
type Token struct {
	Name    string
	Literal string
}

func (t Token) GetLiteral() string {
	return t.Literal
}