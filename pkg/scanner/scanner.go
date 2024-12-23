package scanner

import (
	scanner "asritha.dev/compiler/pkg/scanner/regex" //fix this:https://stackoverflow.com/questions/45899203/can-i-develop-a-go-package-in-multiple-source-directories
	"fmt"
)

type Scanner struct {
	tInitInfos []*TokenInitInfo
}

func (s *Scanner) convertTokenRegexToFA() {
	for _, tII := range s.tInitInfos {
		parseTree := convertRegexToParseTree(tII.Regex)
		parseTree.getNFA() //.removeEpsilonTransitions().convertToPseudoDFA().minimize()
	}
}

// TODO: move these helper methods somewhere else ?
func convertRegexToParseTree(regex string) scanner.RExpr {
	return nil
}

func (s *Scanner) Scan() ([]Token, error) {
	return make([]Token, 0), fmt.Errorf("Failed to scan code :(")
}

func NewScanner(tInitInfos []*TokenInitInfo) *Scanner {
	s := &Scanner{
		tInitInfos: tInitInfos,
	}
	s.convertTokenRegexToFA()
	return s
}
