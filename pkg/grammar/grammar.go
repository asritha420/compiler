package grammar

import "asritha.dev/compiler/pkg/utils"

type grammar struct {
	rules []*rule

	ruleNTMap  map[string][]*rule
	firstSets  map[string]map[symbol]struct{}
	followSets map[string]map[symbol]struct{}
}

func NewGrammar(rules ...*rule) *grammar {
	g := &grammar{
		rules:      rules,
		firstSets:  make(map[string]map[symbol]struct{}),
		followSets: make(map[string]map[symbol]struct{}),
		ruleNTMap:  make(map[string][]*rule),
	}

	g.initializeSets()
	g.generateFirstSets()
	g.generateFollowSets()

	return g
}

func (g *grammar) initializeSets() {
	for _, r := range g.rules {
		if _, ok := g.firstSets[r.nonTerm]; !ok {
			g.firstSets[r.nonTerm] = make(map[symbol]struct{})
			g.followSets[r.nonTerm] = make(map[symbol]struct{})
			g.ruleNTMap[r.nonTerm] = make([]*rule, 0)
		}
		g.ruleNTMap[r.nonTerm] = append(g.ruleNTMap[r.nonTerm], r)
	}
}

func (g *grammar) generateFirstSet(sententialForm ...*symbol) map[symbol]struct{} {
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
			utils.AddToMapIgnore(g.firstSets[symbol.name], firstSet, Epsilon)
			if _, containsEpsilon := g.firstSets[symbol.name][Epsilon]; !containsEpsilon {
				break sententialLoop
			}
			sententialFormIdx++
		}
	}
	return firstSet
}

func (g *grammar) generateFirstSets() {
	changeMade := true
	for changeMade {
		changeMade = false
		for _, rule := range g.rules {
			newFirstSet := g.generateFirstSet(rule.sententialForm...)
			if utils.AddToMap(newFirstSet, g.firstSets[rule.nonTerm]) != 0 {
				changeMade = true
			}
		}
	}
}

func (g *grammar) generateFollowSets() {
	// add EOF to first rule
	g.followSets[g.rules[0].nonTerm][EndOfInput] = struct{}{}

	changeMade := true
	for changeMade {
		changeMade = false
		for _, rule := range g.rules {
			for i, s := range rule.sententialForm {
				if s.symbolType != nonTerm {
					continue
				}

				firstSet := g.generateFirstSet(rule.sententialForm[i+1:]...)
				_, containsEpsilon := firstSet[Epsilon]
				delete(firstSet, Epsilon)
				if utils.AddToMap(firstSet, g.followSets[s.name]) != 0 {
					changeMade = true
				}
				if containsEpsilon && utils.AddToMap(g.followSets[rule.nonTerm], g.followSets[s.name]) != 0 {
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