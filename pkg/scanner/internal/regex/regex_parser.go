package regex

import (
	"asritha.dev/compiler/pkg/grammar"
	"fmt"
)


/* 

	rules = []string{
		`production = expression`,
		`expression = term expressionPrime`,
		`expressionPrime = "|" term expressionPrime | " "`, //TODO: write in the spec how EPSILON should be specified as " "
		`term = factor termPrime`,
		`termPrime = factor termPrime | " "`,
		`factor = group factorPrime`,
		`factorPrime = "*" factorPrime | " "`,
		`group = "(" expression ")" | [a-z] | [A-Z] | [0-9]`,
	}

// parse tree for the regex grammar 
		 
	"abc"

	expr 
	concat var string random 
	a b c
	

	var string random 
	

	for each production: 
	symbolTokens = []symbol { 
		expression nt token, 
	}

	bababa

	*/ 

// RegexParser will consume a regex string and return an AST
type RegexParser struct {
	grammar *grammar.Grammar
	regex   []rune 
	index   int //index points to the current unconsumed byte in regex
}

func NewRegexParser(grammar *grammar.Grammar, regex string) *RegexParser {
	return &RegexParser{grammar: grammar, regex: []rune(regex), index: 0}
}

// consumeIf() returns true and consumes the current byte in regex if it matches any of the specified bytes; otherwise, it returns false and the parser state remains unchanged
func (p *RegexParser) consumeIf(matches ...rune) (*rune, bool) {
	if len(p.regex) <= p.index {
		return nil, false
	}

	curr := p.regex[p.index]
	for _, b := range matches {
		if curr == b {
			p.index++
			return &curr, true
		}
	}
	return nil, false
}

// lookAhead() returns true if current byte in regex matches any of the specified bytes
func (p *RegexParser) lookAhead(matches ...rune) (*rune, bool) {
	if len(p.regex) <= p.index {
		return nil, false
	}

	for _, b := range matches {
		if b == p.regex[p.index] {
			return nil, true
		}
	}
	return nil, false
}

/*
RegexParser performs the following checks:
	1) if rule X cannot produce epsilon, and RegexParser encounters a character NOT in FIRST(X), there is a syntax error
	2) if rule X can produce epsilon, and RegexParser encounters a character NOT in FIRST(X), it accepts rule X -> epsilon
*/

/*
Parse() is entrypoint to parser
P -> E
*/
func (p *RegexParser) Parse() (RExpr, error) {
	return p.parseE()
}

// E -> TE'
func (p *RegexParser) parseE() (RExpr, error) {
	t, err := p.parseT()
	if err != nil {
		return nil, err
	}
	for {
		ePrime, err := p.parseEPrime()
		if err != nil {
			return nil, err
		}
		if ePrime != nil {
			t = NewAlternation(t, ePrime)
			continue
		}
		break
	}

	return t, nil
}

/*
E' -> |TE'
E' -> epsilon
*/
func (p *RegexParser) parseEPrime() (RExpr, error) {
	if _, ok := p.consumeIf(p.grammar.Rules['X'].FirstSet...); ok {
		return p.parseT()
	}
	return nil, nil
}

// T -> FT'
func (p *RegexParser) parseT() (RExpr, error) {
	f, err := p.parseF()
	if err != nil {
		return nil, err
	}

	for {
		tPrime, err := p.parseTPrime()
		if err != nil {
			return nil, err
		}
		if tPrime != nil {
			f = NewConcatenation(f, tPrime)
			continue
		}
		break
	}

	return f, nil
}

/*
T' -> FT'
T' -> epsilon
*/
func (p *RegexParser) parseTPrime() (RExpr, error) {
	if _, ok := p.lookAhead(p.grammar.Rules['Y'].FirstSet...); ok {
		return p.parseF()
	}
	return nil, nil
}

// F -> GF'
func (p *RegexParser) parseF() (RExpr, error) {
	g, err := p.parseG()
	if err != nil {
		return nil, err
	}

	for p.parseFPrime() {
		g = NewKleeneStar(g)
	}

	return g, nil
}

/*
F' -> *F'
F' -> epsilon
*/
func (p *RegexParser) parseFPrime() bool {
	_, ok := p.consumeIf(p.grammar.Rules['R'].FirstSet...)
	return ok
}

/*
G -> (E)
G -> v //where v is any byte
*/
func (p *RegexParser) parseG() (RExpr, error) {
	if _, ok := p.consumeIf('('); ok {
		e, err := p.parseE()
		if err != nil {
			return nil, err
		}
		p.consumeIf(')')
		return e, nil
	}

	if b, ok := p.consumeIf(p.grammar.Rules['G'].FirstSet...); ok { //TODO: fix indexes
		return NewConst(*b), nil
	}

	return nil, fmt.Errorf("parse error: unexpected byte %c", p.regex[p.index])
}
