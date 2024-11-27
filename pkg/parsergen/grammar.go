package parsergen

import (
	"fmt"
)

/*
  P -> E
  E -> TE'
  E' -> "|"TE' | ""
  T -> FT'
  T' -> FT' | "\\\""
  F -> GF'
  F' -> "*"F' | ""
  G -> "("E")" | "v" // where v is any byte
*/

type Grammar struct {
	Rules        []*Rule
	nonTerminals [][]rune
	terminals    []rune
}

func NewGrammar(rulesStr string) (*Grammar, error) {
	rules, nonTerminals, err := NewRules(rulesStr)
	if err != nil {
		return nil, fmt.Errorf("[NewGrammar] failed to generate rules because:\n%w", err)
	}

	g := &Grammar{
		Rules:        rules,
		nonTerminals: nonTerminals,
		terminals:    make([]rune, 0),
	}
	
	g.fillTerminals()

	return g, nil
}

func (g *Grammar) fillTerminals() {
	for _, rule := range g.Rules {
		for _, production := range rule.productions {
			for _, token := range production {
				if token.tokenType == TERMINAL {
					g.terminals = append(g.terminals, token.value...)
				}
			}
		}
	}
}

//
//// generateFirstSet() will recursively calculate the FIRST set for each rule in the Grammar;  FIRST(x) is the set of all terminals (including epsilon) that can appear at the beginning of any possible derivation of rule x.
//func (g *Grammar) generateFirstSet(startRName byte, rName byte) {
//	for _, r := range g.Rules[rName] {
//		//if the rule is an epsilon or a valid character or int, add it to FIRST(rName)
//		if r == Epsilon {
//			g.FirstSet[rName] = append(g.FirstSet[rName], 0x00)
//		} else if r == ValidChar {
//			g.FirstSet[rName] = append(g.FirstSet[rName], ValidChar...)
//			g.FirstSet[startRName] = append(g.FirstSet[startRName], ValidChar...)
//		} else if r == ValidInt {
//			g.FirstSet[rName] = append(g.FirstSet[rName], ValidInt...)
//			g.FirstSet[startRName] = append(g.FirstSet[startRName], ValidInt...)
//		} else if slices.Contains(g.terminals, r[0]) {
//			g.FirstSet[rName] = append(g.FirstSet[rName], r[0]) // if the first byte of Rules[rName] is a terminal, add it to FIRST(rName)
//			g.FirstSet[startRName] = append(g.FirstSet[startRName], r[0])
//		} else if slices.Contains(g.nonTerminals, r[0]) {
//			g.generateFirstSet(startRName, r[0]) // if the first byte of Rules[rName] is a nonTerminal, recursively call generateFirstSet() until it finds a terminal
//		}
//	}
//}
//func (g *Grammar) generateFirstSet(rule *Rule, nextRule *Rule) {
//
//	for _, rule := range g.Rules {
//		for _, production := range rule.Productions {
//			if production == Epsilon {
//				rule.firstSet = append(rule.firstSet, Epsilon)
//			} else if production == ValidInt || production == ValidChar {
//				rule.firstSet = append(rule.firstSet, production)
//			}
//		}
//	}
//}
