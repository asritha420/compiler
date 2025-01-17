package grammar

import (
	"testing"
)

func TestFirstFollow(t *testing.T) {
	E := NewNonTerm("E")
	EP := NewNonTerm("E'")
	T := NewNonTerm("T")
	TP := NewNonTerm("T'")
	F := NewNonTerm("F")

	plus := NewToken("+")
	i := NewToken("int")
	lParen := NewToken("(")
	rParen := NewToken(")")
	mult := NewToken("*")

	r1 := NewRule("P", E)
	r2 := NewRule("E", T, EP)
	r3 := NewRule("E'", plus, T, EP)
	r4 := NewRule("E'", Epsilon)
	r5 := NewRule("T", F, TP)
	r6 := NewRule("T'", mult, F, TP)
	r7 := NewRule("T'", Epsilon)
	r8 := NewRule("F", lParen, E, rParen)
	r9 := NewRule("F", i)

	g := NewGrammar(r1, r2, r3, r4, r5, r6, r7, r8, r9)

	println(g)
}

func setupGrammar() *Grammar {
	E := NewNonTerm("E")
	T := NewNonTerm("T")

	plus := NewToken("+")
	id := NewToken("id")
	lParen := NewToken("(")
	rParen := NewToken(")")

	r1 := NewRule("P", E)
	r2 := NewRule("E", E, plus, T)
	r3 := NewRule("E", T)
	r4 := NewRule("T", id, lParen, E, rParen)
	r5 := NewRule("T", id)

	return NewGrammar(r1, r2, r3, r4, r5)
}

// func TestLR1(t *testing.T) {
// 	g := setupGrammar()
// 	_, states := generateLR1(g)
// 	print(makeMermaid(states))
// }

// func TestLALR(t *testing.T) {
// 	g := setupGrammar()
// 	_, states := generateLALR(g)
// 	print(makeMermaid(states))
// }

// func TestParser(t *testing.T) {
// 	p := NewParser(setupGrammar(), true)
// 	input := []TokenSymbol{
// 		{name: "id"},
// 		{name: "("},
// 		{name: "id"},
// 		{name: "+"},
// 		{name: "id"},
// 		{name: ")"},
// 	}
// 	tree, _ := p.Parse(input)
// 	print(tree)
// }
