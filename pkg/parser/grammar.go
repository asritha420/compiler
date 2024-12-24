package parser

import "slices"

// TODO: go through and determine what should be public and what should not and also organize into its own files
type Rule struct {
	Terminal    string
	NonTerminal string

	nonTerminalRuneSlice []rune
	nonTerminalTokens    []RuleNonTerminalToken
}

type RuleNonTerminalScanner struct {
	start int //first character of grammar symbol being currently considered
	curr  int //unconsumed char being currently considered
	rule  *Rule
	g     *Grammar
}

// Each Token represents a grammar symbol
type RuleNonTerminalToken struct {
	isTerminal bool //true = terminal; false = nonTerminal
	literal    string
}

func NewRule(terminal string, nonTerminal string) *Rule {
	r := &Rule{
		Terminal:             terminal,
		NonTerminal:          nonTerminal,
		nonTerminalRuneSlice: []rune(nonTerminal),
	}
	return r
}

func (rs *RuleNonTerminalScanner) Scan() {
	for !rs.isAtEnd() {

	}
}

func (rs *RuleNonTerminalScanner) ScanRuleToken() {
	c := rs.advance()

	switch c {
	case '"':
		//belongs to a terminal -> keep going until it finds the next ", then verify that it matches something listed in []terminals, else throw an error -> invalid terminal
		while
	default:
		//is the start of a nonTerminal, keep going until it matches the longest possible string listed in the nonTerminals; if it reaches EOF without matching w/ a valid nonTerminal, throw an error

	}
}

// consumeCurrRune() consumes and returns the current character
func (rs *RuleNonTerminalScanner) consumeCurrRune() rune {
	rs.curr++
	return rs.rule.nonTerminalRuneSlice[rs.curr-1]
}

// TODO: add comments for all these
func (rs *RuleNonTerminalScanner) isValidTerminal(gs string) bool {
	return slices.Contains(rs.g.Terminals, gs)
}

func (rs *RuleNonTerminalScanner) isValidNonTerminal(gs string) bool {
	return slices.Contains(rs.g.NonTerminals, gs)
}

func (rs *RuleNonTerminalScanner) advance() rune {
	rs.curr++
	r := rs.rule.nonTerminalRuneSlice[rs.curr-1]
	return r
}

func (rs *RuleNonTerminalScanner) isAtEnd() bool {
	if rs.curr != len(rs.rule.nonTerminalRuneSlice)-1 {
		return false
	}
	return true
}

type Grammar struct {
	Rules        []*Rule
	NonTerminals []string
	Terminals    []string
}

func NewGrammar(rules []*Rule, nonTerminals []string, terminals []string) *Grammar {
	g := &Grammar{
		Rules:        rules,
		NonTerminals: nonTerminals,
		Terminals:    terminals,
	}

	for _, r := range rules {
		rs := RuleNonTerminalScanner{
			start: 0,
			curr:  0,
			rule:  r,
		}
		rs.Scan()
	}

	return g
}

/*
	Given grammar symbol X

	FIRST(X) is the set of all terminal symbols that can appear at the beginning of any string derived from X

	For a terminal X:
		FIRST(X) = {X}
	For a production rule X -> Y_1Y_2...Y_n
		For Y_i in Y_1Y_2...Y_n
			Add all terminal symbols in FIRST(Y_i) to FIRST(X) except epsilon
			If epsilon is in FIRST(Y_i), continue to Y_i+1
		If all Y_1, Y_2, ... Y_n can derive epsilon, add epsilon to FIRST(X)

*/

func (g *Grammar) generateFirstSets() {

	//for i := len(g.Rules) - 1; i >= 0; i-- {
	//	r := g.Rules[i]
	//
	//}
}

func (g *Grammar) generateFollowSets() {
	//TODO
}
