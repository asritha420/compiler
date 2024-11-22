package parsergen

import (
	"slices"
	"strings"
)

type grammar struct {
	Rules        []Rule
	nonTerminals []string
	terminals    []string //TODO: should terminals be []byte ?
}

func newGrammar(rules []Rule) *grammar {
	g := &grammar{
		Rules:        rules,
		terminals:    make([]string, 0),
		nonTerminals: make([]string, 0),
	}
	g.fillNonTerminals()
	g.fillTerminals()
	return g
}

func (g *grammar) fillNonTerminals() {
	for _, rule := range g.Rules {
		g.nonTerminals = append(g.nonTerminals, rule.NonTerminal)
	}
}

func (g *grammar) fillTerminals() {
	for _, rule := range g.Rules {
		for _, production := range rule.Productions {
			//if production == int || production == ValidChar {
			//	g.terminals = append(g.terminals, production...)
			//}

			//TODO: how to check for number and letter, and append it to terminals

			quoteIndex := 0
			for quoteIndex != -1 && quoteIndex < (len(production)-2) {
				quoteIndex = strings.Index(production, "\"")
				if quoteIndex != -1 {
					g.terminals = append(g.terminals, production[quoteIndex+1:quoteIndex+2])
				}
				production = production[quoteIndex+3:]
			}
		}
	}
}

// generateFirstSet() will recursively calculate the FIRST set for each rule in the grammar;  FIRST(x) is the set of all terminals (including epsilon) that can appear at the beginning of any possible derivation of rule x.
func (g *grammar) generateFirstSet(startRName byte, rName byte) {
	for _, r := range g.Rules[rName] {
		//if the rule is an epsilon or a valid character or int, add it to FIRST(rName)
		if r == Epsilon {
			g.FirstSet[rName] = append(g.FirstSet[rName], 0x00)
		} else if r == ValidChar {
			g.FirstSet[rName] = append(g.FirstSet[rName], ValidChar...)
			g.FirstSet[startRName] = append(g.FirstSet[startRName], ValidChar...)
		} else if r == ValidInt {
			g.FirstSet[rName] = append(g.FirstSet[rName], ValidInt...)
			g.FirstSet[startRName] = append(g.FirstSet[startRName], ValidInt...)
		} else if slices.Contains(g.terminals, r[0]) {
			g.FirstSet[rName] = append(g.FirstSet[rName], r[0]) // if the first byte of Rules[rName] is a terminal, add it to FIRST(rName)
			g.FirstSet[startRName] = append(g.FirstSet[startRName], r[0])
		} else if slices.Contains(g.nonTerminals, r[0]) {
			g.generateFirstSet(startRName, r[0]) // if the first byte of Rules[rName] is a nonTerminal, recursively call generateFirstSet() until it finds a terminal
		}
	}
}
func (g *grammar) generateFirstSet(rule *Rule, nextRule *Rule) {

	for _, rule := range g.Rules {
		for _, production := range rule.Productions {
			if production == Epsilon {
				rule.firstSet = append(rule.firstSet, Epsilon)
			} else if production == ValidInt || production == ValidChar {
				rule.firstSet = append(rule.firstSet, production)
			}
		}
	}
}
