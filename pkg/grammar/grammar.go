package grammar

import (
	"fmt"
	"strings"
)

// TODO: users should not have to specify terminals and nonTerminals?
// TODO: test NewGrammar()
// TODO: if its the symbol, store a a pointer to the same symbol in the Rule

type Grammar struct {
	Rules []*Rule
}

func NewGrammar(rules []string, nonTerminals []string, terminals []string) (*Grammar, error) {
	g := new(Grammar)
	rs := RuleScanner{
		validNonTerminals: nonTerminals,
		validTerminals:    terminals,
	}

	for i, ruleString := range rules {
		lhs, rhs, containsAssignmentOperator := strings.Cut(ruleString, "=")
		if !containsAssignmentOperator {
			return nil, fmt.Errorf("NewGrammar(): rule #%d ('%s') does not contain the assignment operator '='", i, ruleString)
		}

		rs.curr = 0
		rs.rule = []rune(strings.TrimSpace(rhs))
		productions, err := rs.Scan()
		if err != nil {
			return nil, fmt.Errorf("NewGrammar(): rule #%d ('%s'): %v", i, ruleString, err)
		}

		nonTerminal := strings.TrimSpace(lhs)
		if !rs.isValidNonTerminal(nonTerminal) {
			return nil, fmt.Errorf("NewGrammar(): rule #%d ('%s') contains invalid non-terminal '%s' on L.H.S", i, ruleString, nonTerminal)
		}

		g.Rules = append(g.Rules, &Rule{
			nonTerminal: nonTerminal,
			productions: productions,
		})
	}

	//g.generateFirstSets()
	//g.generateFollowSets()

	return g, nil
}
