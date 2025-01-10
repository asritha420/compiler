package grammar

import "asritha.dev/compiler/pkg/utils"

type Grammar struct {
	Rules []*rule

	ruleNTMap  map[string][]*rule
	FirstSets  map[string]map[symbol]struct{}
	FollowSets map[string]map[symbol]struct{}
}

func NewGrammar(rules ...*rule) *Grammar {
	g := &Grammar{
		Rules:      rules,
		FirstSets:  make(map[string]map[symbol]struct{}),
		FollowSets: make(map[string]map[symbol]struct{}),
		ruleNTMap:  make(map[string][]*rule),
	}

	g.initializeSets()
	g.generateFirstSets()
	g.generateFollowSets()

	return g
}

func (g *Grammar) initializeSets() {
	for _, r := range g.Rules {
		if _, ok := g.FirstSets[r.nonTerm]; !ok {
			g.FirstSets[r.nonTerm] = make(map[symbol]struct{})
			g.FollowSets[r.nonTerm] = make(map[symbol]struct{})
			g.ruleNTMap[r.nonTerm] = make([]*rule, 0)
		}
		g.ruleNTMap[r.nonTerm] = append(g.ruleNTMap[r.nonTerm], r)
	}
}

func (g *Grammar) generateFirstSet(sententialForm ...*symbol) map[symbol]struct{} {
	firstSet := make(map[symbol]struct{})
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
			if utils.AddToMap(newFirstSet, g.FirstSets[rule.nonTerm]) != 0 {
				changeMade = true
			}
		}
	}
}

func (g *Grammar) generateFollowSets() {
	// add EOF to first rule
	g.FollowSets[g.Rules[0].nonTerm][EndOfInput] = struct{}{}

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
				if containsEpsilon && utils.AddToMap(g.FollowSets[rule.nonTerm], g.FollowSets[s.name]) != 0 {
					changeMade = true
				}
			}
		}
	}
}

// func (g *grammar) canProduceEpsilon(sententialForm ...*symbol) bool {
// 	for _, s := range sententialForm {
// 		switch s.symbolType {
// 		case endOfInput, token:
// 			return false
// 		case nonTerm:
// 			if _, containsEpsilon := g.firstSets[s.name][Epsilon]; !containsEpsilon {
// 				return false
// 			}
// 		}
// 	}
// 	return true
// }
