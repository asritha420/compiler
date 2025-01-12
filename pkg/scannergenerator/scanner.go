package scannergenerator

import (
	"regexp"
)

// for now all, must have an identifer token for the bug
type Scanner struct {
	tokenSpec []TokenInfo // in order of token priority
}

type TokenInfo struct {
	TokenType   string
	regexp      regexp.Regexp
	RegexString string
}

func (s *Scanner) Scan(code string) ([]Token, error) {

	tokenStream := make([]Token, 0)
	currToken := &Token{}
	var currWord string

	matchesNone := func() bool {
		for _, ts := range s.tokenSpec {
			if ts.regexp.MatchString(currWord) {
				currToken.Name = ts.TokenType
				currToken.Literal = currWord
				return false
			}
		}
		return true
	}

	for _, character := range code {
		currWord += string(character)
		if matchesNone() {
			tokenStream = append(tokenStream, *currToken)
			currWord = string(character)
		}
	}

	return tokenStream, nil
}

func NewScanner(tokenSpec []TokenInfo) (*Scanner, error) {
	s := &Scanner{
		tokenSpec: make([]TokenInfo, 0),
	}

	for _, tokenInfo := range s.tokenSpec {
		regex, err := regexp.Compile(tokenInfo.RegexString)
		if err != nil {
			return nil, err
		}
		tokenInfo.regexp = *regex
	}

	return s, nil
}
