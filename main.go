package main

import (
	"fmt"

	"asritha.dev/compiler/pkg/parsergen"
)

func main() {

	// // (b | ab | aab)*
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
	// // DFA, _ := scannergen.ConvertNFAtoDFA(NFA)
	// // scannergen.minimizeDFA(DFA, DFAMap)
	// println(scannergen.MakeMermaid(NFA))

	test, err := parsergen.ConvertProductions("ERTE'T\" \\\"7\\\\\" | E'ER TE'", []string{"ER", "ERT", "T", "E'", "E", "RE😊"})
	if err != nil {
		fmt.Print(err)
		return
	}
	// \"7\\
	for _, t := range test {
		fmt.Printf("%c", t)
	}
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
