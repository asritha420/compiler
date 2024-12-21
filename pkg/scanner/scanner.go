package scanner

import (
	"fmt"
)

type Scanner struct {
	tInitInfos []*TokenInitInfo
}

func (s *Scanner) convertTokenRegexToFA() {
	for _, tII := range s.tInitInfos {

	}
}

// TODO: move these helper methods somewhere else ?
func convertRegexToParseTree(regex string) *RExpr {

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
