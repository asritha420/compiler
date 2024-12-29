package scanner

import (
	"asritha.dev/compiler/pkg/scanner/internal/regex"
	"fmt"
)

type Scanner struct {
	string
	Tokens []Token
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

// Scan() should have support for both code file and REPL
func (s *Scanner) Scan(code string) ([]Token, error) {
	return make([]Token, 0), fmt.Errorf("Failed to scan code :(")
}

func NewScanner(tInitInfos []*TokenInitInfo) *Scanner {
	s := &Scanner{
		//tInitInfos: tInitInfos,
	}
	s.convertTokenRegexToFA()
	return s
}
