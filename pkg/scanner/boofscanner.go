package scanner

//import (
//	"regexp"
//)

//type TokenType int

//
//// for now all, must have an identifer token for the bug
//type ScannerGen struct {
//	tokenSpec []TokenInfo // in order of token priority
//}
//
//type TokenInfo struct {
//	TokenType   string
//	regexp      regexp.Regexp
//	RegexString string
//}
//
//func (s *Scanner) Scan(code string) ([]Token, error) {
//
//	tokenStream := make([]Token, 0)
//	currToken := &Token{}
//	currWord := ""
//	nextIdx := 0
//
//	findMatch := func(str string) bool {
//		for _, ts := range s.tokenSpec {
//			if len(ts.regexp.FindString(str)) == len(str) {
//				currToken.Name = ts.TokenType
//				currToken.Literal = str
//				return true
//			}
//		}
//		return false
//	}
//
//	for nextIdx != len(code) {
//		currWord += string(code[nextIdx])
//		nextIdx++
//		if !findMatch(currWord) {
//			if len(currWord) == 1 {
//				tokenStream = append(tokenStream, Token{currWord, currWord})
//			} else {
//				tokenStream = append(tokenStream, *currToken)
//				nextIdx--
//			}
//			currWord = ""
//		}
//	}
//
//	if len(currWord) != 0 && findMatch(currWord) {
//		tokenStream = append(tokenStream, *currToken)
//	}
//
//	return tokenStream, nil
//}
//
//func NewScanner(tokenSpec []TokenInfo) (*Scanner, error) {
//	s := &Scanner{
//		tokenSpec: make([]TokenInfo, 0),
//	}
//
//	for _, tokenInfo := range tokenSpec {
//		regex, err := regexp.Compile(tokenInfo.RegexString)
//		if err != nil {
//			return nil, err
//		}
//		tokenInfo.regexp = *regex
//		s.tokenSpec = append(s.tokenSpec, tokenInfo)
//	}
//
//	return s, nil
//}
