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
