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
	rules []*rule

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
	for _, rule := range g.rules {
		if _, ok := g.firstSets[rule.nonTerm]; !ok {
			g.firstSets[rule.nonTerm] = make(map[symbol]struct{})
			g.followSets[rule.nonTerm] = make(map[symbol]struct{})
		}
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

func main() {
	E := newSymbol(nonTerm, "E")
	EP := newSymbol(nonTerm, "E'")
	T := newSymbol(nonTerm, "T")
	TP := newSymbol(nonTerm, "T'")
	F := newSymbol(nonTerm, "F")

	plus := newSymbol(terminal, "+")
	i := newSymbol(token, "int")
	lParen := newSymbol(terminal, "(")
	rParen := newSymbol(terminal, ")")
	mult := newSymbol(terminal, "*")

	r1 := newRule("P", E)
	r2 := newRule("E", T, EP)
	r3 := newRule("E'", plus, T, EP)
	r4 := newRule("E'", &epsilon)
	r5 := newRule("T", F, TP)
	r6 := newRule("T'", mult, F, TP)
	r7 := newRule("T'", &epsilon)
	r8 := newRule("F", lParen, E, rParen)
	r9 := newRule("F", i)

	g := newGrammar(r1, r2, r3, r4, r5, r6, r7, r8, r9)

	print(g)
}
