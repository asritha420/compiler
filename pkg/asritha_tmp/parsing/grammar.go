package aparsing

import "slices"

type Rule struct {
	nonTerminal rune
	production  []rune
}

func NewRule(nonTerminal rune, production string) *Rule {
	return &Rule{nonTerminal, []rune(production)}
}

type LRGrammar struct {
	rules                   []*Rule
	nonTerminals, terminals []rune
}

func (g *LRGrammar) isValidNonTerminal(symbol rune) bool {
	return slices.Contains(g.nonTerminals, symbol)
}

func (g *LRGrammar) isValidTerminal(symbol rune) bool {
	return slices.Contains(g.terminals, symbol)
}

func (g *LRGrammar) getRulesForNT(nonTerminal rune) []*Rule {
	matchingRules := make([]*Rule, 0)
	for _, rule := range g.rules {
		if rule.nonTerminal == nonTerminal {
			matchingRules = append(matchingRules, rule)
		}
	}
	return matchingRules
}
