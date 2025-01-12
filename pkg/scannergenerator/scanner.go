package scannergenerator

import (
	"regexp"
	"strings"
)

type Scanner struct {
	tokenSpec map[TokenType]regexp.Regexp
}

func (s *Scanner) Scan(code string) ([]Token, error) {

	tokenStream := make([]Token, 0)

	codeWords := strings.Split(code, " ")

	for _, word := range codeWords {
		for tokenType, regex := range s.tokenSpec {
			if regex.MatchString(word) {
				tokenStream = append(tokenStream, Token{
					TokenType: tokenType,
					Lexeme:    word,
				})
			}
		}
	}

	return tokenStream, nil
}

func NewScanner(tokenSpec map[TokenType]string) (*Scanner, error) {
	s := &Scanner{
		tokenSpec: make(map[TokenType]regexp.Regexp),
	}
	for tokenType, regexString := range tokenSpec {
		regex, err := regexp.Compile(regexString)
		if err != nil {
			return nil, err
		}
		s.tokenSpec[tokenType] = *regex
	}
	return s, nil
}
