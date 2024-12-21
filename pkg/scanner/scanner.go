package scanner

import "fmt"

type Scanner struct {
	tInitInfos []*TokenInitInfo
}

func (s *Scanner) convertTokenRegexToFA() {

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
