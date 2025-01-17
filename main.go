package main

// TODO fix all the pointers!!

import (
	"fmt"

	. "asritha.dev/compiler/pkg/grammar"
	. "asritha.dev/compiler/pkg/parser"
	"asritha.dev/compiler/pkg/scanner"
)

// func GenerateGrammarScanner() *Scanner {

// 	grammarTokens := []TokenInfo{
// 		{
// 			TokenType:   "letter",
// 			RegexString: "[a-zA-Z]",
// 		},
// 		{
// 			TokenType:   "digit",
// 			RegexString: "[0-9]",
// 		},
// 		{
// 			TokenType:   "space",
// 			RegexString: `\s+`,
// 		},
// 	}
// 	grammarScanner, _ := NewScanner(grammarTokens)

// 	return grammarScanner
// }

func GenerateGrammar() *Grammar {
	//scanner tokens
	letter := NewToken("letter")
	digit := NewToken("digit")
	space := NewToken("space")

	//other tokens
	leftBracket := NewToken("[")
	rightBracket := NewToken("]")
	leftBracketCurly := NewToken("{")
	rightBracketCurly := NewToken("}")
	leftParen := NewToken("(")
	rightParen := NewToken(")")
	lessThan := NewToken("<")
	greaterThan := NewToken(">")
	singleQuote := NewToken("'")
	doubleQuote := NewToken("\"")
	equal := NewToken("=")
	pipe := NewToken("|")
	dot := NewToken(".")
	comma := NewToken(",")
	semicolon := NewToken(";")
	dash := NewToken("-")
	plus := NewToken("+")
	asterisk := NewToken("*")
	question := NewToken("?")
	underscore := NewToken("_")
	forwardSlash := NewToken("/")
	backSlash := NewToken("\\")

	//non terms
	symbol := NewNonTerm("symbol")
	strChar := NewNonTerm("strChar")
	idChar := NewNonTerm("idChar")
	identifier := NewNonTerm("identifier")
	str := NewNonTerm("string")
	token := NewNonTerm("token")
	separator := NewNonTerm("separator")
	term := NewNonTerm("term")
	sTerm := NewNonTerm("sTerm")
	unary := NewNonTerm("unary")
	factor := NewNonTerm("factor")
	concatenation := NewNonTerm("concatenation")
	alternation := NewNonTerm("alternation")
	lhs := NewNonTerm("lhs")
	rhs := NewNonTerm("rhs")
	rule := NewNonTerm("rule")
	rules := NewNonTerm("rules")
	tokenCap := NewNonTerm("tokenCap")

	//rules
	//symbol
	s1 := NewRule("symbol", leftBracket)
	s2 := NewRule("symbol", rightBracket)
	s3 := NewRule("symbol", leftBracketCurly)
	s4 := NewRule("symbol", rightBracketCurly)
	s5 := NewRule("symbol", leftParen)
	s6 := NewRule("symbol", rightParen)
	s7 := NewRule("symbol", lessThan)
	s8 := NewRule("symbol", greaterThan)
	s9 := NewRule("symbol", singleQuote)
	s10 := NewRule("symbol", backSlash, doubleQuote)
	s11 := NewRule("symbol", equal)
	s12 := NewRule("symbol", pipe)
	s13 := NewRule("symbol", dot)
	s14 := NewRule("symbol", comma)
	s15 := NewRule("symbol", semicolon)
	s16 := NewRule("symbol", dash)
	s17 := NewRule("symbol", plus)
	s18 := NewRule("symbol", asterisk)
	s19 := NewRule("symbol", question)
	s20 := NewRule("symbol", underscore)
	s21 := NewRule("symbol", forwardSlash)
	s22 := NewRule("symbol", backSlash, backSlash)

	//strChar
	sc1 := NewRule("strChar", letter)
	sc2 := NewRule("strChar", digit)
	sc3 := NewRule("strChar", space)
	sc4 := NewRule("strChar", symbol)

	//idChar
	ic1 := NewRule("idChar", letter)
	ic2 := NewRule("idChar", digit)
	ic3 := NewRule("idChar", underscore)

	//id
	id1 := NewRule("identifier", idChar)
	id2 := NewRule("identifier", idChar, identifier)

	//string
	str1 := NewRule("string", Epsilon)
	str2 := NewRule("string", strChar, str)

	//token
	tc1 := NewRule("tokenCap", doubleQuote)
	tok1 := NewRule("token", tokenCap, str, tokenCap)

	//separator
	se1 := NewRule("separator", space)
	se2 := NewRule("separator", Epsilon)

	//term
	t1 := NewRule("term", token)
	t2 := NewRule("term", identifier)
	t3 := NewRule("term", leftParen, rhs, rightParen)

	//sTerm
	st1 := NewRule("sTerm", separator, term, separator)

	//unary
	u1 := NewRule("unary", question)
	u2 := NewRule("unary", asterisk)
	u3 := NewRule("unary", plus)
	u4 := NewRule("unary", Epsilon)

	//factor
	f1 := NewRule("factor", sTerm, unary, separator)

	//concatenation
	c1 := NewRule("concatenation", factor)
	c2 := NewRule("concatenation", factor, comma, concatenation)

	//alternation
	a1 := NewRule("alternation", concatenation)
	a2 := NewRule("alternation", concatenation, pipe, alternation)

	//other
	lhs1 := NewRule("lhs", identifier)
	rhs1 :=  NewRule("rhs", alternation)

	r1 := NewRule("rule", separator, lhs, separator, equal, rhs, semicolon, separator)
	rs1 := NewRule("rules", rule)
	rs2 := NewRule("rules", rule, rules)

	start := NewRule("start", rules)


	return NewGrammar(
		start,
		s1,s2,s3,s4,s5,s6,s7,s8,s9,s10,s11,s12,s13,s14,s15,s16,s17,s18,s19,s20,s21,s22,
		sc1,sc2,sc3,sc4,
		ic1,ic2,ic3,
		id1,id2,
		str1,str2,
		tc1, tok1,
		se1,se2,
		t1,t2,t3,
		st1,
		u1,u2,u3,u4,
		f1,
		c1,c2,
		a1,a2,
		lhs1,rhs1,
		r1,rs1,rs2,
	)
}

type grammarAST interface {
}

type grammarRule struct {
	nt string
	prod grammarAST
}

func main() {
	g := GenerateGrammar()

	// gs := *GenerateGrammarScanner()

	p := NewParser(g, true)

	// tokens, err := gs.Scan("P=E ;\n\nE=(\"te\\\"st\",lol)*;\nR=hello | test, y\n;")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	tokens := []scanner.Token{
		{Name: "letter", Literal: "P"},
		{Name: "=", Literal: "="},
		{Name: "letter", Literal: "E"},
		{Name: ";", Literal: ";"},
	}
	tree, _ := p.Parse(tokens)
	
	fmt.Println(tree.GetLiteral())
}
