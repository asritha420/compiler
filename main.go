package main

import (
	. "asritha.dev/compiler/pkg/grammar"
	. "asritha.dev/compiler/pkg/parser"
	"fmt"
	"log"
)

func main() {
	g := GenerateGrammar()

	gs := GetGrammarScanner()

	p := NewParser(g, true)

	tokens, err := gs.Scan("P = E; \n E = (lol);")
	if err != nil {
		log.Fatal(err)
	}
	p.Parse(tokens)
	fmt.Println(tokens)
}
