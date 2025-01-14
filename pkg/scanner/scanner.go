package scanner

import (
	"asritha.dev/compiler/pkg/regex"
)

type Scanner interface {
	Scan()
}

type ScannerError struct {
	message string
	offset  int
}

func (e *ScannerError) Error() string {
	return e.message
}

// need to be list, not map
func New(map[string]regex.Regex) Scanner {
	return nil
}

type scannerGen struct {
}
