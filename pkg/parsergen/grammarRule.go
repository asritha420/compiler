package parsergen

import (
	"fmt"
	"slices"
)

type RuleTokenType int

const (
	TERMINAL RuleTokenType = iota
	NON_TERMINAL
)

type Rule struct {
	NonTerminal string
	Productions [][]RuleToken
	firstSet    []string
	followSet   []string
	terminals   []rune
}

type RuleToken struct {
	tokenType RuleTokenType
	value []rune
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

func ConvertProductions(rule string, validNTs []string) ([][]RuleToken, error) {
	ruleRunes := []rune(rule)
	validNTsRunes := make([][]rune, 0)
	for _, validNT := range validNTs{
		validNTsRunes = append(validNTsRunes, []rune(validNT))
	} 
	
	productions := make([][]RuleToken, 0)
	production := make([]RuleToken, 0)

	isTerminal := false
	isEscaped := false //Should only be true if in terminal mode!!

	ntStart := -1
	var ntString []rune 
	var currValidNTs [][]rune

	for i, ch := range ruleRunes {
		if isEscaped {
			if ch != '"' && ch != '\\' {
				return nil, fmt.Errorf("invalid character '%c' after \\ at index %d", ch, i)
			}

			production = append(production, RuleToken{tokenType: TERMINAL, value: []rune{ch}})
			isEscaped = false
			continue
		}
		if ch == '"' {
			isTerminal = !isTerminal
			continue
		}
		if isTerminal {
			if ch == '\\' {
				isEscaped = true
				continue
			}
			production = append(production, RuleToken{tokenType: TERMINAL, value: []rune{ch}})
			continue
		}

		// skip non terminal spaces (only if not looking for non terminal)
		if ch == ' ' {
			if ntStart != -1 {
				return nil, fmt.Errorf("space encountered at index %d without finishing non terminal", i)
			}
			continue
		}

		// check for new production
		if ch == '|' {
			if ntStart != -1 {
				// invalid state
				return nil, fmt.Errorf("start of new production at index %d without finishing previous one", i)
			}
			productions = append(productions, production)
			production = make([]RuleToken, 0)
			continue
		}

		// start new NT
		if ntStart == -1 {
			ntStart = i
			ntString = ruleRunes[ntStart:]
			currValidNTs = nonTerminalLookAhead(ntString, 0, validNTsRunes) // filters first character
		}

		// look ahead 1
		lookAheadIndex := i - ntStart + 1
		tmpValidNTs := nonTerminalLookAhead(ntString, lookAheadIndex, currValidNTs)
		if len(tmpValidNTs) == 0 {
			// NT no longer valid
			ntString = ntString[:lookAheadIndex]
			if !nonTerminalMatches(ntString, currValidNTs){
				return nil, fmt.Errorf("invalid non-terminal \"%s\" at index %d", string(ntString), i)
			}

			production = append(production, RuleToken{tokenType: NON_TERMINAL, value: ntString})
			ntStart = -1
			continue
		}
		currValidNTs = tmpValidNTs
	}
	
	if(isTerminal){
		return nil, fmt.Errorf("end of rule reached without closing \"")
	}

	productions = append(productions, production)

	return productions, nil
}

/*
NTLookAhead will take the string p and check if the character at an index matches any NT at the same index.
*/
func nonTerminalLookAhead(p []rune, index int, allNTs [][]rune) [][]rune {
	validNTs := make([][]rune, 0)
	if len(p) <= index{
		return validNTs
	}

	for _, nt := range allNTs {
		if len(nt) > index && nt[index] == p[index] {
			validNTs = append(validNTs, nt)
		}
	}
	return validNTs
}

func nonTerminalMatches(nt []rune, validNTs [][]rune) bool {
	for _, validNT := range validNTs{
		if slices.Equal(nt, validNT) {
			return true
		}
	}
	return false
}