package scanner

import "testing"

func TestNewScanner(t *testing.T) {
	// Tests w/ Tiny C
	const (
		// Keywords
		FOR TokenType = iota
		WHILE
		DO
		IF
		ELSE
	)

	_ = []TokenInitInfo{
		{
			TType: FOR,
			Name:  "for",
			Regex: "for",
		},
	}

}

func TestScan(t *testing.T) {

}
