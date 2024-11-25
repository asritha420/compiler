package parsergen

import "fmt"

type Rule struct {
	NonTerminal string
	Productions [][]RuleToken
	firstSet    []string
	followSet   []string
	terminals   []rune
}

type RuleToken interface {
}

type Terminal struct {
	value rune
}

type NonTerminal struct {
	value string
}

//
//func NewRule(nT string, productions []string) *Rule {
//	return &Rule{
//		NonTerminal: nT,
//		Productions: productions,
//	}
//}

//func NewRules(rules string) []*Rule {
/*
	for each line

	lh of arrow: becomes nt string
	rh of arrow: is prodcution, will convert into tokens

*/
//}

func ConvertProduction(p string, validNT []string) ([]RuleToken, error) {
	rules := make([]RuleToken, 0)
	isTerminal := false
	ntStart := -1
	var ntString string
	isEscaped := false
	for i, ch := range p {
		if ch == '"' {
			if isEscaped && isTerminal {
				isEscaped = false
			} else {
				isTerminal = !isTerminal
				continue
			}
		}
		if isTerminal {
			if ch == '\\' {
				if isEscaped {
					isEscaped = false
				} else {
					isEscaped = true
					continue
				}
			}
			rules = append(rules, Terminal{value: ch})
		} else {
			if ntStart == -1 {
				ntStart = i
				ntString = p[ntStart:]
			}

			if !prodLookAhead(ntString, i-ntStart+1, validNT) {
				// found new NT
				ntString = ntString[:i-ntStart+1]
				if err := matchesNT(ntString, validNT); err != nil {
					return nil, err
				}
				rules = append(rules, NonTerminal{value: ntString})
				ntStart = -1
			}
		}
	}
	return rules, nil
}

func prodLookAhead(p string, index int, validNT []string) bool {
	for i := 0; i < len(validNT); i++ {
		if len(validNT[i]) > index && validNT[i][index] == p[index] {
			return true
		}
	}
	return false
}

func matchesNT(p string, validNT []string) error {
	matched := false
	for _, nt := range validNT {
		if p == nt {
			if matched {
				return fmt.Errorf("multiple tokens matched %s", p)
			}
			matched = true
		}
	}
	if !matched {
		return fmt.Errorf("token %s not found", p)
	} else {
		return nil
	}
}
