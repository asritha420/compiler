package grammar

import (
	"fmt"
)

type grammar struct {
	rules     []*rule
	ruleNTMap map[string][]*rule

	firstSets  map[string]map[symbol]struct{}
	followSets map[string]map[symbol]struct{}
}

func newGrammar(rules ...*rule) *grammar {
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

func addToSet(src, dst map[symbol]struct{}) bool {
	added := false
	for s := range src {
		if _, ok := dst[s]; !ok {
			dst[s] = struct{}{}
			added = true
		}
	}
	return added
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
		switch symbol.sType {
		case Terminal, Token:
			firstSet[*symbol] = struct{}{}
			// weird edge case where they put more symbols after an epsilon symbol
			if *symbol == Epsilon {
				sententialFormIdx++
			} else {
				break sententialLoop
			}

		case NonTerm:
			containsEpsilon := false
			for s := range g.firstSets[symbol.data] {
				if s == Epsilon {
					containsEpsilon = true
				}
				firstSet[s] = struct{}{}
			}
			if !containsEpsilon {
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
			if addToSet(newFirstSet, g.firstSets[rule.nonTerm]) {
				changeMade = true
			}
		}
	}
}

func (g *grammar) generateFollowSets() {
	// add EOF to first rule
	g.followSets[g.rules[0].nonTerm][EndOfFile] = struct{}{}

	changeMade := true
	for changeMade {
		changeMade = false
		for _, rule := range g.rules {
			for i, s := range rule.sententialForm {
				if s.sType != NonTerm {
					continue
				}

				firstSet := g.generateFirstSet(rule.sententialForm[i+1:]...)
				_, containsEpsilon := firstSet[Epsilon]
				delete(firstSet, Epsilon)
				if addToSet(firstSet, g.followSets[s.data]) {
					changeMade = true
				}
				if containsEpsilon && addToSet(g.followSets[rule.nonTerm], g.followSets[s.data]) {
					changeMade = true
				}
			}
		}
	}
}

type augmentedRule struct {
	rule      *rule
	position  int
	lookahead map[symbol]struct{}
}

func newAugmentedRule(r *rule, position int) *augmentedRule {
	return &augmentedRule{
		rule:      r,
		position:  position,
		lookahead: make(map[symbol]struct{}),
	}
}

/*
Returns a copy of the augment rule with the position shifted one to the right (next symbol)
*/
func (ar *augmentedRule) shiftedCopy() *augmentedRule {
	if ar.position == len(ar.rule.sententialForm) {
		return nil
	}
	return newAugmentedRule(ar.rule, ar.position+1)
}

/*
Returns the next symbol (symbol to right of position) in an augmented rule or nil if there is no next symbol
*/
func (ar *augmentedRule) getNextSymbol() *symbol {
	//Note: position should NOT be more than len(sententialForm)
	if ar.position == len(ar.rule.sententialForm) {
		return nil
	}
	return ar.rule.sententialForm[ar.position]
}

func (ar augmentedRule) String() string {
	// TODO add lookahead
	output := ar.rule.nonTerm + "="
	for i, s := range ar.rule.sententialForm {
		if ar.position == i {
			output += "."
		}
		output += s.String()
	}
	if ar.position == len(ar.rule.sententialForm) {
		output += "."
	}
	return output
}

type lr1AutomationState struct {
	id             uint
	augmentedRules map[*augmentedRule]struct{}
	transitions    map[symbol]*lr1AutomationState
}

func newLR1AutomationState(id *uint) *lr1AutomationState {
	*id++
	return &lr1AutomationState{
		id:             *id - 1,
		augmentedRules: make(map[*augmentedRule]struct{}),
		transitions:    make(map[symbol]*lr1AutomationState),
	}
}

func (g *grammar) getClosureRecursion(ar *augmentedRule, closure map[*augmentedRule]struct{}, closed map[string]struct{}) {
	if _, ok := closed[ar.String()]; ok {
		return
	}

	closed[ar.String()] = struct{}{}
	closure[ar] = struct{}{}

	nextSymbol := ar.getNextSymbol()
	if nextSymbol == nil || nextSymbol.sType != NonTerm {
		return
	}

	// nextSymbol is a NT
	for _, r := range g.ruleNTMap[nextSymbol.data] {
		newAR := newAugmentedRule(r, 0)
		g.getClosureRecursion(newAR, closure, closed)
	}
}

func (g *grammar) getClosure(ars ...*augmentedRule) map[*augmentedRule]struct{} {
	closure := make(map[*augmentedRule]struct{})
	closed := make(map[string]struct{})

	for _, ar := range ars {
		g.getClosureRecursion(ar, closure, closed)
	}

	return closure
}

func (g *grammar) getTransitions(ars map[*augmentedRule]struct{}) map[symbol]map[*augmentedRule]struct{} {
	transitions := make(map[symbol]map[*augmentedRule]struct{})
	closed := make(map[symbol]map[string]struct{})
	for ar := range ars {
		nextSymbol := ar.getNextSymbol()
		if nextSymbol == nil {
			continue
		}

		if _, ok := transitions[*nextSymbol]; !ok {
			transitions[*nextSymbol] = make(map[*augmentedRule]struct{})
			closed[*nextSymbol] = make(map[string]struct{})
		}
		g.getClosureRecursion(ar.shiftedCopy(), transitions[*nextSymbol], closed[*nextSymbol])
	}
	return transitions
}

func equal(m1, m2 map[*augmentedRule]struct{}) bool {
	if len(m1) != len(m2) {
		return false
	}

	strs := make(map[string]struct{})
	for ar := range m1 {
		strs[ar.String()] = struct{}{}
	}

	for ar := range m2 {
		if _, ok := strs[ar.String()]; !ok {
			return false
		}
	}

	return true
}

func findState(target map[*augmentedRule]struct{}, states []*lr1AutomationState) *lr1AutomationState {
	for _, state := range states {
		fmt.Printf("target: %v\n", target)
		fmt.Printf("state : %v\n", state.augmentedRules)
		if equal(target, state.augmentedRules) {
			return state
		}
	}
	return nil
}

func (g *grammar) generateLR1() *lr1AutomationState {
	var id uint = 0

	kernel := newLR1AutomationState(&id)
	startRule := newAugmentedRule(g.rules[0], 0)
	kernel.augmentedRules = g.getClosure(startRule)

	states := []*lr1AutomationState{kernel}
	openList := []*lr1AutomationState{kernel}

	for len(openList) > 0 {
		state := openList[0]
		openList = openList[1:]

		for transition, rules := range g.getTransitions(state.augmentedRules) {
			if transitionState := findState(rules, states); transitionState != nil {
				state.transitions[transition] = transitionState
				continue
			}

			newState := newLR1AutomationState(&id)
			newState.augmentedRules = rules
			state.transitions[transition] = newState
			states = append(states, newState)
			openList = append(openList, newState)
		}
	}

	return kernel
}

func main() {
	// E := newSymbol(nonTerm, "E")
	// EP := newSymbol(nonTerm, "E'")
	// T := newSymbol(nonTerm, "T")
	// TP := newSymbol(nonTerm, "T'")
	// F := newSymbol(nonTerm, "F")

	// plus := newSymbol(terminal, "+")
	// i := newSymbol(token, "int")
	// lParen := newSymbol(terminal, "(")
	// rParen := newSymbol(terminal, ")")
	// mult := newSymbol(terminal, "*")

	// r1 := newRule("P", E)
	// r2 := newRule("E", T, EP)
	// r3 := newRule("E'", plus, T, EP)
	// r4 := newRule("E'", &epsilon)
	// r5 := newRule("T", F, TP)
	// r6 := newRule("T'", mult, F, TP)
	// r7 := newRule("T'", &epsilon)
	// r8 := newRule("F", lParen, E, rParen)
	// r9 := newRule("F", i)

	// g := newGrammar(r1, r2, r3, r4, r5, r6, r7, r8, r9)

	// print(g)

	E := newSymbol(NonTerm, "E")
	T := newSymbol(NonTerm, "T")

	plus := newSymbol(Terminal, "+")
	id := newSymbol(Token, "id")
	lParen := newSymbol(Terminal, "(")
	rParen := newSymbol(Terminal, ")")

	r1 := newRule("P", E)
	r2 := newRule("E", E, plus, T)
	r3 := newRule("E", T)
	r4 := newRule("T", id, lParen, E, rParen)
	r5 := newRule("T", id)

	g := newGrammar(r1, r2, r3, r4, r5)
	graph := g.generateLR1()
	fmt.Print(graph)
}
