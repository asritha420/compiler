package grammar

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



// func main() {
// 	// E := newSymbol(nonTerm, "E")
// 	// EP := newSymbol(nonTerm, "E'")
// 	// T := newSymbol(nonTerm, "T")
// 	// TP := newSymbol(nonTerm, "T'")
// 	// F := newSymbol(nonTerm, "F")

// 	// plus := newSymbol(terminal, "+")
// 	// i := newSymbol(token, "int")
// 	// lParen := newSymbol(terminal, "(")
// 	// rParen := newSymbol(terminal, ")")
// 	// mult := newSymbol(terminal, "*")

// 	// r1 := newRule("P", E)
// 	// r2 := newRule("E", T, EP)
// 	// r3 := newRule("E'", plus, T, EP)
// 	// r4 := newRule("E'", &epsilon)
// 	// r5 := newRule("T", F, TP)
// 	// r6 := newRule("T'", mult, F, TP)
// 	// r7 := newRule("T'", &epsilon)
// 	// r8 := newRule("F", lParen, E, rParen)
// 	// r9 := newRule("F", i)

// 	// g := newGrammar(r1, r2, r3, r4, r5, r6, r7, r8, r9)

// 	// print(g)

// 	E := newSymbol(NonTerm, "E")
// 	T := newSymbol(NonTerm, "T")

// 	plus := newSymbol(Terminal, "+")
// 	id := newSymbol(Token, "id")
// 	lParen := newSymbol(Terminal, "(")
// 	rParen := newSymbol(Terminal, ")")

// 	r1 := newRule("P", E)
// 	r2 := newRule("E", E, plus, T)
// 	r3 := newRule("E", T)
// 	r4 := newRule("T", id, lParen, E, rParen)
// 	r5 := newRule("T", id)

// 	g := newGrammar(r1, r2, r3, r4, r5)
// 	graph := g.generateLR1()
// 	fmt.Print(graph)
// }
