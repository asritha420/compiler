package grammar

import (
	"log"

	"asritha.dev/compiler/pkg/scannergenerator"
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

	return g
}

func NewSimpleGrammar(rules ...string) *Grammar {
	g := &Grammar{
		FirstSets:  make(map[string]utils.Set[Symbol]),
		FollowSets: make(map[string]utils.Set[Symbol]),
	}

	return g
}

func GetGrammarScanner() *scannergenerator.Scanner {

	grammarTokens := []scannergenerator.TokenInfo{
		{
			TokenType:   "letter",
			RegexString: "[a-zA-Z]",
		},
		{
			TokenType:   "digit",
			RegexString: "[0-9]",
		},
		{
			TokenType:   "space",
			RegexString: `(" " | "\n" | "\t" | "\r" | "\f" | "\b")*`,
		},
	}
	grammarScanner, err := scannergenerator.NewScanner(grammarTokens)
	if err != nil {
		log.Fatal(err)
	}
	
	return grammarScanner
}

func GenerateGrammar() *Grammar {
	// Non-terminals
	rules := NewNonTerm("rules")
	rule := NewNonTerm("rule")
	lhs := NewNonTerm("lhs")
	rhs := NewNonTerm("rhs")
	alternation := NewNonTerm("alternation")
	concatenation := NewNonTerm("concatenation")
	term := NewNonTerm("term")
	terminal := NewNonTerm("terminal")
	chars := NewNonTerm("chars")
	character := NewNonTerm("character")
	identifier := NewNonTerm("identifier")
	identifierChar := NewNonTerm("identifierChar")

	// Tokens
	assign := NewToken("=")
	semicolon := NewToken(";")
	pipe := NewToken("|")
	comma := NewToken(",")
	epsilon := &Epsilon

	// Symbol tokens
	symbol := NewToken("symbol")
	letter := NewToken("letter")
	digit := NewToken("digit")
	underscore := NewToken("_")
	space := NewToken(" ")

	// Rules
	r1 := NewRule("start", rules)
	r2 := NewRule("rules", rule, rules)
	r3 := NewRule("rules", rule)
	r4 := NewRule("rule", lhs, assign, rhs, semicolon)
	r5 := NewRule("lhs", identifier)
	r6 := NewRule("rhs", alternation)
	r7 := NewRule("alternation", concatenation, pipe, alternation)
	r8 := NewRule("alternation", concatenation)
	r9 := NewRule("concatenation", term, comma, concatenation)
	r10 := NewRule("concatenation", term)
	r11 := NewRule("term", terminal)
	r12 := NewRule("term", identifier)
	r13 := NewRule("term", NewToken("("), rhs, NewToken(")"))
	r14 := NewRule("terminal", NewToken("\""), chars, NewToken("\""))
	r15 := NewRule("chars", character, chars)
	r16 := NewRule("chars", epsilon)
	r17 := NewRule("character", letter)
	r18 := NewRule("character", digit)
	r19 := NewRule("character", symbol)
	r20 := NewRule("character", underscore)
	r21 := NewRule("character", space)
	r22 := NewRule("identifier", identifierChar, identifier)
	r23 := NewRule("identifier", identifierChar)
	r24 := NewRule("identifierChar", letter)
	r25 := NewRule("identifierChar", digit)
	r26 := NewRule("identifierChar", underscore)

	// Create grammar
	g := NewGrammar(r1, r2, r3, r4, r5, r6, r7, r8, r9, r10, r11, r12, r13, r14, r15, r16, r17, r18, r19, r20, r21, r22, r23, r24, r25, r26)
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
