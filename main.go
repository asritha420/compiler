package main

import (
	"asritha.dev/compiler/pkg/parsergen"
	"asritha.dev/compiler/pkg/scannergen"
	"fmt"
	"log"
)

func main() {
	TestRegexParser()
}
func TestRegexParser() {
	regexGrammar := parsergen.NewLL1Grammar(
		map[byte][]string{
			'P': {"E"},
			'E': {"TX"},
			'X': {"|TX", parsergen.Epsilon},
			'T': {"FY"},
			'Y': {"FY", parsergen.Epsilon},
			'F': {"GM"},
			'M': {"*M", parsergen.Epsilon},
			'G': {"(E)", parsergen.ValidChar, parsergen.ValidInt},
		},
	)
	regexGrammar.FollowSet = map[byte][]byte{'t': []byte{'t'}} //manually overwriting for now as generateFollowSet() is not completed

	regexTests := []string{"a(b|c)*", "cat|cow", "a(b|a|0|d)*a", "a**", "(abb)*", "a(cow|cat)*"}

	for _, rt := range regexTests {
		rParser := scannergen.NewRegexParser(regexGrammar, rt)
		ast, err := rParser.Parse()
		if err != nil {
			fmt.Printf("❌ PARSING ERROR: %v\n", err)
		}
		fmt.Printf("%+v: \n", ast)

		if val, ok := ast.(scannergen.ASTPrinter); ok {
			fmt.Printf("%v \n \n", val.PrintNode(""))
		}
	}

	testRParser := scannergen.NewRegexParser(regexGrammar, "a(cow|cat)*")
	ast, err := testRParser.Parse()
	if err != nil {
		fmt.Printf("❌ PARSING ERROR: %v\n", err)
	}

	graph, _, _, err := scannergen.ConvertRegexToNfa(ast)
	if err != nil {
		log.Fatal(err)
	}
	DFA, _ := scannergen.ConvertNFAtoDFA(graph)
	println(scannergen.MakeMermaid(DFA))
}
