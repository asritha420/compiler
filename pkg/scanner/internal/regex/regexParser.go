package regex

// //TODO: a lot of this stuff shouldnt be public. same with the fa package. go through the entire scanner package and remove stuff that shouldnt be public

// import (
// 	"fmt"
//     "asritha.dev/compiler/pkg/parser"
// )

// // RegexParser will consume a regex string and return an AST
// type RegexParser struct {
// 	grammar *parser.LL1Grammar
// 	regex   string
// 	index   int //index points to the current unconsumed byte in regex
// }

// func NewRegexParser(grammar *parsergen.LL1Grammar, regex string) *RegexParser {
// 	return &RegexParser{grammar: grammar, regex: regex, index: 0}
// }

// // consumeIf() returns true and consumes the current byte in regex if it matches any of the specified bytes; otherwise, it returns false and the parser state remains unchanged
// func (p *RegexParser) consumeIf(matches ...byte) (*byte, bool) {
// 	if len(p.regex) <= p.index {
// 		return nil, false
// 	}

// 	curr := p.regex[p.index]
// 	for _, b := range matches {
// 		if curr == b {
// 			p.index++
// 			return &curr, true
// 		}
// 	}
// 	return nil, false
// }

// // lookAhead() returns true if current byte in regex matches any of the specified bytes
// func (p *RegexParser) lookAhead(matches ...byte) (*byte, bool) {
// 	if len(p.regex) <= p.index {
// 		return nil, false
// 	}

// 	for _, b := range matches {
// 		if b == p.regex[p.index] {
// 			return nil, true
// 		}
// 	}
// 	return nil, false
// }

// /*
// RegexParser performs the following checks:
// 	1) if rule X cannot produce epsilon, and RegexParser encounters a character NOT in FIRST(X), there is a syntax error
// 	2) if rule X can produce epsilon, and RegexParser encounters a character NOT in FIRST(X), it accepts rule X -> epsilon
// */

// /*
// Parse() is entrypoint to parser
// P -> E
// */
// func (p *RegexParser) Parse() (RExpr, error) {
// 	return p.parseE()
// }

// // E -> TE'
// func (p *RegexParser) parseE() (RExpr, error) {
// 	t, err := p.parseT()
// 	if err != nil {
// 		return nil, err
// 	}
// 	//printmermaidNFA()
// }
// 	for {
// 		ePrime, err := p.parseEPrime()
// 		if err != nil {
// 			return nil, err
// 		}
// 		if ePrime != nil {
// 			t = NewAlternation(t, ePrime)
// 			continue
// 		}
// 		break
// 	}

// 	return t, nil
// }

// /*
// E' -> |TE'
// E' -> epsilon
// */
// func (p *RegexParser) parseEPrime() (RExpr, error) {
// 	if _, ok := p.consumeIf(p.grammar.FirstSet['X']...); ok {
// 		return p.parseT()
// 	}
// 	return nil, nil
// }

// // T -> FT'
// func (p *RegexParser) parseT() (RExpr, error) {
// 	f, err := p.parseF()
// 	if err != nil {
// 		return nil, err
// 	}

// 	for {
// 		tPrime, err := p.parseTPrime()
// 		if err != nil {
// 			return nil, err
// 		}
// 		if tPrime != nil {
// 			f = NewConcatenation(f, tPrime)
// 			continue
// 		}
// 		break
// 	}

// 	return f, nil
// }

// /*
// T' -> FT'
// T' -> epsilon
// */
// func (p *RegexParser) parseTPrime() (RExpr, error) {
// 	if _, ok := p.lookAhead(p.grammar.FirstSet['Y']...); ok {
// 		return p.parseF()
// 	}
// 	return nil, nil
// }

// // F -> GF'
// func (p *RegexParser) parseF() (RExpr, error) {
// 	g, err := p.parseG()
// 	if err != nil {
// 		return nil, err
// 	}

// 	for p.parseFPrime() {
// 		g = NewKleeneStar(g)
// 	}

// 	return g, nil
// }

// /*
// F' -> *F'
// F' -> epsilon
// */
// func (p *RegexParser) parseFPrime() bool {
// 	_, ok := p.consumeIf(p.grammar.FirstSet['M']...)
// 	return ok
// }

// /*
// G -> (E)
// G -> v //where v is any byte
// */
// func (p *RegexParser) parseG() (RExpr, error) {
// 	if _, ok := p.consumeIf('('); ok {
// 		e, err := p.parseE()
// 		if err != nil {
// 			return nil, err
// 		}
// 		p.consumeIf(')')
// 		return e, nil
// 	}

// 	if b, ok := p.consumeIf(p.grammar.FirstSet['G']...); ok {
// 		return NewConst(*b), nil
// 	}

// 	return nil, fmt.Errorf("parse error: unexpected byte %c", p.regex[p.index])
// }
