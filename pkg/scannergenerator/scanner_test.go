package scannergenerator

import (
	"slices"
	"testing"
)

// TODO: how am I handling errors thrown
// TODO: limit length of identifier
// TODO: figure out what should be public and private

func TestNewScanner(t *testing.T) {
	const (
		// Keywords
		FOR TokenType = iota
		IF
		IDENTIFIER
		ADD
		NUMBER
	)

	correctTokens := []Token{
		{
			TokenType: FOR,
			Lexeme:    "for",
			Line:      0,
		},
		{
			TokenType: IF,
			Lexeme:    "if",
			Line:      0,
		},
		{
			TokenType: NUMBER,
			Lexeme:    "234234",
			Line:      0,
		},
		{
			TokenType: ADD,
			Lexeme:    "+",
			Line:      0,
		},
		{
			TokenType: IDENTIFIER,
			Lexeme:    "asdf_123",
			Line:      0,
		},
	}

	tokensSpec := map[TokenType]string{
		FOR:        "for",
		IF:         "if",
		IDENTIFIER: "[a-zA-Z0-9_-]+",
		ADD:        "+",
		NUMBER:     "[0-9]+",
	}

	// TODO: move to test file
	// TODO: check edge cases in this
	testCode := `
		for if 234234 + asdf_123
	`

	scanner, err := NewScanner(tokensSpec)
	if err != nil {

	}
	tokens, err := scanner.Scan(testCode)
	if err != nil {
		t.Errorf("Scan(): invalid tokens. %w", err)
	}

	if slices.Equal(tokens, correctTokens) {
		t.Errorf("Scan(): invalid tokens. %w", tokens) // TODO: fix this
	}
}

func TestScan(t *testing.T) {
}
