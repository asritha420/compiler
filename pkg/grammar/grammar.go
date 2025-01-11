package grammar

import "asritha.dev/compiler/pkg/utils"

type Grammar struct {
	Rules []*rule

	firstRule *rule
	ruleNTMap  map[string][]*rule
	FirstSets  map[string]set[symbol]
	FollowSets map[string]set[symbol]
}

func NewGrammar(rules ...*rule) *Grammar {
	g := &Grammar{
		Rules:      rules,
		firstRule: rules[0],
		FirstSets:  make(map[string]set[symbol]),
		FollowSets: make(map[string]set[symbol]),
		ruleNTMap:  make(map[string][]*rule),
	}

	g.initializeSets()
	g.generateFirstSets()
	g.generateFollowSets()

	return g
}

func (g *Grammar) initializeSets() {
	for _, r := range g.Rules {
		if _, ok := g.FirstSets[r.NonTerm]; !ok {
			g.FirstSets[r.NonTerm] = make(set[symbol])
			g.FollowSets[r.NonTerm] = make(set[symbol])
			g.ruleNTMap[r.NonTerm] = make([]*rule, 0)
		}
		g.ruleNTMap[r.NonTerm] = append(g.ruleNTMap[r.NonTerm], r)
	}
}

func (g *Grammar) generateFirstSet(sententialForm ...*symbol) set[symbol] {
	firstSet := make(set[symbol])
	sententialFormIdx := 0
sententialLoop:
	for {
		if sententialFormIdx == len(sententialForm) {
			firstSet[Epsilon] = struct{}{}
			break sententialLoop
		}

		symbol := sententialForm[sententialFormIdx]
		switch symbol.symbolType {
		case epsilon:
			sententialFormIdx++

		// This really shouldn't be inside a rule (it should only be used in follow sets)
		case endOfInput:
			break sententialLoop

		case token:
			firstSet[*symbol] = struct{}{}
			break sententialLoop

		case nonTerm:
			utils.AddToMapIgnore(g.FirstSets[symbol.name], firstSet, Epsilon)
			if _, containsEpsilon := g.FirstSets[symbol.name][Epsilon]; !containsEpsilon {
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
			newFirstSet := g.generateFirstSet(rule.sententialForm...)
			if utils.AddToMap(newFirstSet, g.FirstSets[rule.NonTerm]) != 0 {
				changeMade = true
			}
		}
	}
}

func (g *Grammar) generateFollowSets() {
	// add EOF to first rule
	g.FollowSets[g.Rules[0].NonTerm][EndOfInput] = struct{}{}

	changeMade := true
	for changeMade {
		changeMade = false
		for _, rule := range g.Rules {
			for i, s := range rule.sententialForm {
				if s.symbolType != nonTerm {
					continue
				}

				firstSet := g.generateFirstSet(rule.sententialForm[i+1:]...)
				_, containsEpsilon := firstSet[Epsilon]
				delete(firstSet, Epsilon)
				if utils.AddToMap(firstSet, g.FollowSets[s.name]) != 0 {
					changeMade = true
				}
				if containsEpsilon && utils.AddToMap(g.FollowSets[rule.NonTerm], g.FollowSets[s.name]) != 0 {
					changeMade = true
				}
			}
		}
	}
}