package parsergen

import (
	"fmt"
	"slices"
	"strings"
)

type RuleTokenType int

type RuneRange struct {
}

const (
	TERMINAL RuleTokenType = iota
	NON_TERMINAL
)

type Rule struct {
	nonTerminal string
	productions [][]RuleToken
	firstSet    []string
	followSet   []string
}

type RuleToken struct {
	tokenType RuleTokenType
	value     []rune
}

// function for all constraints on non terminals
// eg. no spaces or quotes
func validNonTerminal(nt string) bool {
	return !strings.ContainsAny(nt, " \n\t\"[]")
}

func NewRules(inputRulesStr string) ([]*Rule, [][]rune, error) {
	inputRules := strings.Split(inputRulesStr, "\n")
	parsedRules := make(map[string]string)

	// split the rule
	for i, rule := range inputRules {
		rule = strings.TrimSpace(rule)
		if len(rule) == 0 {
			continue
		}

		tokens := strings.Split(rule, "->")
		if len(tokens) != 2 {
			return nil, nil, fmt.Errorf("[newRules] invalid rule format at line %d", i+1)
		}
		nt := strings.TrimSpace(tokens[0])
		if !validNonTerminal(nt) {
			return nil, nil, fmt.Errorf("[newRules] non terminal %s at line %d is invalid", nt, i+1)
		}
		if _, ok := parsedRules[nt]; ok {
			return nil, nil, fmt.Errorf("[newRules] repeated non terminal %s at line %d", nt, i+1)
		}

		parsedRules[nt] = strings.TrimSpace(tokens[1])
	}

	// make valid non terminals
	validNTs := make([][]rune, 0)
	for nt, _ := range parsedRules {
		validNTs = append(validNTs, []rune(nt))
	}

	// convert to rules
	rules := make([]*Rule, 0)
	for nt, productions := range parsedRules {
		paresdProductions, err := convertProductions([]rune(productions), validNTs)
		if err != nil {
			return nil, nil, fmt.Errorf("[newRules] failed to convert production %s because:\n%w", productions, err)
		}
		rules = append(rules, &Rule{
			nonTerminal: nt,
			productions: paresdProductions,
		})
	}

	return rules, validNTs, nil
}

func convertProductions(inputProductions []rune, validNTs [][]rune) ([][]RuleToken, error) {
	productions := make([][]RuleToken, 0)
	production := make([]RuleToken, 0)

	isRange := false
	isTerminal := false
	isEscaped := false //Should only be true if in terminal mode!!

	currTokenStart := -1

	var currValidNTs [][]rune

	for i, ch := range inputProductions {
		if isRange {
			// check for end of range
			if ch != ']' {
				continue
			}
			isRange = false
			currRange := inputProductions[currTokenStart:i]
			convertRange(currRange)
		}
		if isEscaped {
			if ch != '"' && ch != '\\' {
				return nil, fmt.Errorf("[convertProductions] invalid character '%c' after \\ at index %d", ch, i)
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

		// check for start of range
		if ch == '[' {
			isRange = true
			currTokenStart = i
			continue
		}

		// skip non terminal spaces (only if not looking for non terminal)
		if ch == ' ' {
			if currTokenStart != -1 {
				return nil, fmt.Errorf("[convertProductions] space encountered at index %d without finishing non terminal", i)
			}
			continue
		}

		// check for new production
		if ch == '|' {
			if currTokenStart != -1 {
				// invalid state
				return nil, fmt.Errorf("[convertProductions] start of new production at index %d without finishing previous one", i)
			}
			productions = append(productions, production)
			production = make([]RuleToken, 0)
			continue
		}

		// start new NT
		if currTokenStart == -1 {
			currTokenStart = i
			currValidNTs = nonTerminalLookAhead(inputProductions[currTokenStart:], 0, validNTs) // filters first character
		}

		// look ahead 1
		lookAheadIndex := i - currTokenStart + 1
		tmpValidNTs := nonTerminalLookAhead(inputProductions[currTokenStart:], lookAheadIndex, currValidNTs)
		if len(tmpValidNTs) == 0 {
			// NT no longer valid
			currNT := inputProductions[currTokenStart:lookAheadIndex]
			if !nonTerminalMatches(currNT, currValidNTs) {
				return nil, fmt.Errorf("[convertProductions] invalid non-terminal %s at index %d", string(currNT), i)
			}

			production = append(production, RuleToken{tokenType: NON_TERMINAL, value: currNT})
			currTokenStart = -1
			continue
		}
		currValidNTs = tmpValidNTs
	}

	if isTerminal {
		return nil, fmt.Errorf("[convertProductions] end of rule reached without closing \"")
	}

	if currTokenStart != -1 {
		return nil, fmt.Errorf("[convertProductions] end of rule reached without closing ]")
	}

	productions = append(productions, production)

	return productions, nil
}

func convertRange(inputRange []rune) ([]RuneRange, error) {
	// languageRange := RuneRange{0, RUNE_MAX}
	return nil, nil
}

/*
NTLookAhead will take the string p and check if the character at an index matches any NT at the same index.
*/
func nonTerminalLookAhead(p []rune, index int, allNTs [][]rune) [][]rune {
	validNTs := make([][]rune, 0)
	if len(p) <= index {
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
	for _, validNT := range validNTs {
		if slices.Equal(nt, validNT) {
			return true
		}
	}
	return false
}
