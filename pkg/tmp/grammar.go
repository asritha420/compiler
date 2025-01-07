package main

type symbolType int

// Note: even if a token is empty (epsilon), it must be passed to the parser.
// Note: end of file should also be passed as a token matching the endOfFile var
const (
	nonTerm symbolType = iota
	terminal
	token
)

var (
	epsilon   = symbol{sType: terminal, data: ""}
	endOfFile = symbol{sType: token, data: "EOF"}
)

/*
Represents a single symbol (can be either a non-terminal or a terminal/symbol)
*/
type symbol struct {
	sType symbolType
	data  string
}

func newSymbol(sType symbolType, data string) *symbol {
	return &symbol{
		sType: sType,
		data:  data,
	}
}

type rule struct {
	nonTerm        string
	sententialForm []*symbol
}

func newRule(nonTerm string, sententialForm ...*symbol) *rule {
	return &rule{
		nonTerm:        nonTerm,
		sententialForm: sententialForm,
	}
}

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
	for {
		if sententialFormIdx == len(sententialForm) {
			firstSet[epsilon] = struct{}{}
			goto sententialEnd
		}

		symbol := sententialForm[sententialFormIdx]
		switch symbol.sType {
		case terminal, token:
			firstSet[*symbol] = struct{}{}
			// weird edge case where they put more symbols after an epsilon symbol
			if *symbol == epsilon {
				sententialFormIdx++
			} else {
				goto sententialEnd
			}

		case nonTerm:
			containsEpsilon := false
			for s := range g.firstSets[symbol.data] {
				if s == epsilon {
					containsEpsilon = true
				}
				firstSet[s] = struct{}{}
			}
			if !containsEpsilon {
				goto sententialEnd
			}

			sententialFormIdx++
		}
	}
sententialEnd:
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
	g.followSets[g.rules[0].nonTerm][endOfFile] = struct{}{}

	changeMade := true
	for changeMade {
		changeMade = false
		for _, rule := range g.rules {
			for i, s := range rule.sententialForm {
				if s.sType != nonTerm {
					continue
				}

				firstSet := g.generateFirstSet(rule.sententialForm[i+1:]...)
				_, containsEpsilon := firstSet[epsilon]
				delete(firstSet, epsilon)
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

func (g *grammar) getNextSymbol(ar *augmentedRule) *symbol {
	//Note: position should NOT be more than len(sententialForm)
	if ar.position == len(ar.rule.sententialForm) {
		if ar.rule == g.rules[0] {
			return &endOfFile
		}
		return nil
	}
	return ar.rule.sententialForm[ar.position]
}

func newAugmentedRule(r *rule, position int) *augmentedRule {
	return &augmentedRule{
		rule:      r,
		position:  position,
		lookahead: make(map[symbol]struct{}),
	}
}

type lr1AutomationState struct {
	id             uint
	augmentedRules []*augmentedRule
	transitions    map[symbol]*lr1AutomationState
}

func newLR1AutomationState(id *uint) *lr1AutomationState {
	*id++
	return &lr1AutomationState{
		id:             *id - 1,
		augmentedRules: make([]*augmentedRule, 0),
		transitions:    make(map[symbol]*lr1AutomationState),
	}
}

func (g *grammar) getClosureRecursion(ar *augmentedRule, closure []*augmentedRule) []*augmentedRule {
	closure = append(closure, ar)

	nextSymbol := g.getNextSymbol(ar)
	if nextSymbol == nil || nextSymbol.sType != nonTerm {
		return closure
	}

	// nextSymbol is a NT
	for _, r := range g.ruleNTMap[nextSymbol.data] {
		newAR := newAugmentedRule(r, 0)
		closure = g.getClosureRecursion(newAR, closure)
	}

	return closure
}

func (g *grammar) getClosure(ars ...*augmentedRule) []*augmentedRule {
	closure := make([]*augmentedRule, 0)
	for _, ar := range ars {
		closure = g.getClosureRecursion(ar, closure)
	}
	return closure
}

func (g *grammar) getNextStates()

func (g *grammar) generateLR1() *lr1AutomationState {
	var id uint = 0
	kernel := newLR1AutomationState(&id)
	startRule := newAugmentedRule(g.rules[0], 0)
	kernel.augmentedRules = g.getClosure(startRule)

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

	E := newSymbol(nonTerm, "E")
	T := newSymbol(nonTerm, "T")

	plus := newSymbol(terminal, "+")
	id := newSymbol(token, "id")
	lParen := newSymbol(terminal, "(")
	rParen := newSymbol(terminal, ")")

	r1 := newRule("P", E)
	r2 := newRule("E", E, plus, T)
	r3 := newRule("E", T)
	r4 := newRule("T", id, lParen, E, rParen)
	r5 := newRule("T", id)

	g := newGrammar(r1, r2, r3, r4, r5)

	print(g)
}
