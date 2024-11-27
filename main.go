package main

import (
	// "asritha.dev/compiler/pkg/scannergen"
	"fmt"

	"asritha.dev/compiler/pkg/parsergen"
)

func main() {

	// (b | ab | aab)*
	// graph := &scannergen.KleeneStar{
	// 	Left: &scannergen.Alternation{
	// 		Left: &scannergen.Const{Value: 'b'},
	// 		Right: &scannergen.Alternation{
	// 			Left: &scannergen.Concatenation{
	// 				Left: &scannergen.Const{Value: 'a'},
	// 				Right: &scannergen.Const{Value: 'b'},
	// 			},
	// 			Right: &scannergen.Concatenation{
	// 				Left: &scannergen.Const{Value: 'a'},
	// 				Right: &scannergen.Concatenation{
	// 					Left: &scannergen.Const{Value: 'a'},
	// 					Right: &scannergen.Const{Value: 'b'},
	// 				},
	// 			},
	// 		},
	// 	},
	// }

	// NFA, _, _, _ := scannergen.ConvertRegexToNfa(graph)
	// DFA, DFAMap := scannergen.ConvertNFAtoDFA(NFA)
	// minDFA := scannergen.MinimizeDFA(DFA.GetId(), DFAMap)
	// println(scannergen.MakeMermaid(minDFA))

	// test, _, err := parsergen.NewRules(`

	// R -> RT
	// T -> RR'RT
	// R' -> "hello" | T'
	// T' -> "world"
	// RT -> "\"" | [helo]

	// `)
	// if err != nil {
	// 	fmt.Print(err)
	// 	return
	// }
	// for _, t := range test {
	// 	fmt.Print(t)
	// }

	ranges := parsergen.MakeRangesThatIgnore(0, 100, 'a','f','9', parsergen.RUNE_MAX-1)
	fmt.Print(ranges)

}

//func RegexParser() {
//	production := parsergen.NewRule("P", []string{"E"})
//	expression := parsergen.NewRule("E", []string{"TE'"})
//	expressionPrime := parsergen.NewRule("E'", []string{"\"|\"TE'", ""})
//	term := parsergen.NewRule("T", []string{"FT'"})
//	termPrime := parsergen.NewRule("T'", []string{"FT'", ""})
//	factor := parsergen.NewRule("F", []string{"GF'"})
//	factorPrime := parsergen.NewRule("F'", []string{"\"*\"F'", ""})
//	group := parsergen.NewRule("G", []string{"\"(\"E\")\"", "v"})
//
//	parsergen.NewGrammar([]*parsergen.Rule{production, expression, expressionPrime, term, termPrime, factor, factorPrime, group})
//}
