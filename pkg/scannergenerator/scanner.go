package scannergenerator

import (
	"regexp"
)

type Scanner struct {
	tokenSpec []TokenInfo // in order of token priority
}

type TokenInfo struct {
	TokenType
	regexp      regexp.Regexp
	regexString string
}

func (s *Scanner) Scan(code string) ([]Token, error) {

	tokenStream := make([]Token, 0)
	currToken := &Token{}
	var currWord string

	matchesNone := func(currWord string) bool {
		for _, ts := range s.tokenSpec {
			if ts.regexp.MatchString(currWord) {
				currToken.TokenType = ts.TokenType
				currToken.Lexeme = currWord
				return false
			}
		}
		return true
	}

	for _, character := range code {
		currWord += string(character)
		if matchesNone(currWord) {
			tokenStream = append(tokenStream, *currToken)
			currWord = ""
		}
	}

	return tokenStream, nil
}

func NewScanner(tokenSpec []TokenInfo) (*Scanner, error) {
	s := &Scanner{
		tokenSpec: make([]TokenInfo, 0),
	}

	for _, tokenInfo := range s.tokenSpec {
		regex, err := regexp.Compile(tokenInfo.regexString)
		if err != nil {
			return nil, err
		}
		tokenInfo.regexp = *regex
	}

	return s, nil
}
