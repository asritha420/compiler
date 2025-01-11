package scannergenerator

import (
	//"asritha.dev/compiler/pkg/scannergenerator/internal/regex"
	"asritha.dev/compiler/pkg/scannergenerator/internal/regex"
	"fmt"
	"regexp"
)

type Scanner struct {
	Tokens    []Token
	tokenSpec map[TokenType]regexp.Regexp
}

func (s *Scanner) convertTokenRegexToFA() {
	//for _, tII := range s.tInitInfos {
	//	_ = convertRegexToParseTree(tII.Regex)
	//	// parseTree.getNFA() //.removeEpsilonTransitions().convertToPseudoDFA().minimize()
	//}
}

// TODO: move these helper methods somewhere else ?
func convertRegexToParseTree(regex string) regex.RExpr {
	return nil
}

func (s *Scanner) Scan(codeFile string) ([]Token, error) {
	return make([]Token, 0), fmt.Errorf("Failed to scan code :(")
}

// NewScanner will take in a map of
func NewScanner(tokenSpecs map[TokenType]string) (*Scanner, error) {
	scanner := &Scanner{
		Tokens:    make([]Token, 0),
		tokenSpec: make(map[TokenType]regexp.Regexp),
	}

	for tokenType, regexString := range tokenSpecs {
		// throw error if regexString is invalid
		regex, err := regexp.Compile(regexString)

		if err != nil {
			return nil, err
		}

		// TODO: check if in our list of supported regexs

		// TODO: handle error of providing multiple definitions for same tokenType

		scanner.tokenSpec[tokenType] = *regex
	}

	return scanner, nil
}
