package grammar

import (
	"asritha.dev/compiler/pkg/utils"
)

type Grammar struct {
	Rules []*Rule

	FirstRule  *Rule
	RuleNTMap  map[string][]*Rule
	FirstSets  map[string]utils.Set[Symbol]
	FollowSets map[string]utils.Set[Symbol]
}

func NewGrammar(rules ...*Rule) *Grammar {
	g := &Grammar{
		Rules:      rules,
		FirstRule:  rules[0],
		FirstSets:  make(map[string]utils.Set[Symbol]),
		FollowSets: make(map[string]utils.Set[Symbol]),
		RuleNTMap:  make(map[string][]*Rule),
	}

	g.initializeSets()
	g.generateFirstSets()
	g.generateFollowSets()

	for _, r := range g.Rules {
		r.removeEpsilon()
	}

	return g
}

func (g *Grammar) initializeSets() {
	for _, r := range g.Rules {
		if _, ok := g.FirstSets[r.NonTerm]; !ok {
			g.FirstSets[r.NonTerm] = make(utils.Set[Symbol])
			g.FollowSets[r.NonTerm] = make(utils.Set[Symbol])
			g.RuleNTMap[r.NonTerm] = make([]*Rule, 0)
		}
		g.RuleNTMap[r.NonTerm] = append(g.RuleNTMap[r.NonTerm], r)
	}
}

func (g *Grammar) GenerateFirstSet(sententialForm ...*Symbol) utils.Set[Symbol] {
	firstSet := make(utils.Set[Symbol])
	sententialFormIdx := 0
sententialLoop:
	for {
		if sententialFormIdx == len(sententialForm) {
			firstSet[Epsilon] = struct{}{}
			break sententialLoop
		}

		symbol := sententialForm[sententialFormIdx]
		switch symbol.SymbolType {
		case epsilonSymbol:
			sententialFormIdx++

		// This really shouldn't be inside a rule (it should only be used in follow sets)
		case endOfInputSymbol:
			break sententialLoop

		case TokenSymbol:
			firstSet[*symbol] = struct{}{}
			break sententialLoop

		case NonTermSymbol:
			utils.AddToMapIgnore(g.FirstSets[symbol.Name], firstSet, Epsilon)
			if _, containsEpsilon := g.FirstSets[symbol.Name][Epsilon]; !containsEpsilon {
				break sententialLoop
			}
			sententialFormIdx++
		}
	}
	return firstSet
}

func (g *Grammar) generateFirstSets() {
	changeMade := true
	for changeMade {
		changeMade = false
		for _, rule := range g.Rules {
			newFirstSet := g.GenerateFirstSet(rule.SententialForm...)
			if utils.AddToMap(newFirstSet, g.FirstSets[rule.NonTerm]) != 0 {
				changeMade = true
			}
		}
	}
}

func (g *Grammar) generateFollowSets() {
	// add EOF to first rule
	g.FollowSets[g.FirstRule.NonTerm][EndOfInput] = struct{}{}

	changeMade := true
	for changeMade {
		changeMade = false
		for _, rule := range g.Rules {
			for i, s := range rule.SententialForm {
				if s.SymbolType != NonTermSymbol {
					continue
				}

				firstSet := g.GenerateFirstSet(rule.SententialForm[i+1:]...)
				_, containsEpsilon := firstSet[Epsilon]
				delete(firstSet, Epsilon)
				if utils.AddToMap(firstSet, g.FollowSets[s.Name]) != 0 {
					changeMade = true
				}
				if containsEpsilon && utils.AddToMap(g.FollowSets[rule.NonTerm], g.FollowSets[s.Name]) != 0 {
					changeMade = true
				}
			}
		}
	}
}
