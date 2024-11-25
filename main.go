package main

import (
	"asritha.dev/compiler/pkg/parsergen"
	"fmt"
)

func main() {

	test, _ := parsergen.ConvertProduction("EE'T\" \\\"7\\\\\"", []string{"E", "T", "E'"})

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
