package main

import (
	"asritha.dev/compiler/pkg/parsergen"
	"asritha.dev/compiler/tests"
	"fmt"
)

func main() {
	//	TestFirstAndFollowSets()
	//	TestRegexParser()

	digits := "01234556789"
	_ = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	test := make([]string, 0)
	for _, x := range digits {
		test = append(test, fmt.Sprintf("\\\"%d\\\"", x))
	}

	fmt.Println(test)
}
func TestFirstAndFollowSets() {
	g9 := tests.TestG9

	/*
		80 = P
		69 = E
		88 = X (E')
		84 = T
		85 = U (T')
		70 = F
	*/

	//TODO: is it ok if these are maps, or will this mess up some order in the code?
	for nt, set := range g9.FirstSet {
		fmt.Printf("FIRST SET %c: ", nt)
		for _, terminal := range set {
			fmt.Printf("%c", terminal)
		}
		fmt.Println()
	}

	for nt, set := range g9.FollowSet {
		fmt.Printf("FOLLOW SET %c: ", nt)
		for _, terminal := range set {
			fmt.Printf("%c", terminal)
		}
		fmt.Println()
	}
}

func TestRegexParser() {
	rules := []parsergen.Rule{
		parsergen.Rule{
			NonTerminal: "P",
			Productions: []string{},
		},
	}
	test := parsergen.NewLL1Grammar()
}

//func TestRegexParser() {
//	//TODO: this should be using the initializer function
//	regexGrammar := parsergen.NewLL1Grammar(
//		map[byte][]string{
//			'P': {"E"},
//			'E': {"TX"},
//			'X': {"\"|\"TX", parsergen.Epsilon},
//			'T': {"FY"},
//			'Y': {"FY", parsergen.Epsilon},
//			'F': {"GM"},
//			'M': {"*M", parsergen.Epsilon},
//			'G': {"(E)", parsergen.ValidChar, parsergen.ValidInt},
//		},
//		[]byte{'|', '*', '(', ')'},
//		[]byte{'P', 'E', 'X', 'T', 'Y', 'F', 'M', 'G'},
//	)
//
//	regexTests := []string{"a(B|c)*", "cat|cow", "a(b|a|0|d)*a", "a**", "(abb)*", "a(cow|cat)*"}
//
//	for _, rt := range regexTests {
//		rParser := scannergen.NewRegexParser(regexGrammar, rt)
//		ast, err := rParser.Parse()
//		if err != nil {
//			fmt.Printf("❌ PARSING ERROR: %v\n", err)
//		}
//		fmt.Printf("%+v: \n", ast)
//
//		if val, ok := ast.(scannergen.ASTPrinter); ok {
//			fmt.Printf("%v \n \n", val.PrintNode(""))
//		}
//	}
//
//	testRParser := scannergen.NewRegexParser(regexGrammar, "a(cow|cat)*")
//	ast, err := testRParser.Parse()
//	if err != nil {
//		fmt.Printf("❌ PARSING ERROR: %v\n", err)
//	}
//
//	graph, _, _, err := scannergen.ConvertRegexToNfa(ast)
//	if err != nil {
//		log.Fatal(err)
//	}
//	DFA, _ := scannergen.ConvertNFAtoDFA(graph)
//	println(scannergen.MakeMermaid(DFA))
//}
