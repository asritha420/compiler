package tests

import "asritha.dev/compiler/pkg/parsergen"

//TODO: make this into an actual testing package later

func tests() {
	//TODO: this should use the Initializer function, also expose as a constant, and change main.go to use this?
	//parsergen.NewLL1Grammar(
	//	map[byte][]string{
	//		'P': {"E"},
	//		'E': {"TX"},
	//		'X': {"|TX", parsergen.Epsilon},
	//		'T': {"FY"},
	//		'Y': {"FY", parsergen.Epsilon},
	//		'F': {"GM"},
	//		'M': {"*M", parsergen.Epsilon},
	//		'G': {"(E)", parsergen.ValidLowercaseChar, parsergen.ValidInt},
	//	},
	//)
	//
}

var (

	/*
	   P -> E
	   E -> TX
	   X -> +TX
	   X -> epsilon
	   T -> FU
	   U -> *FU
	   U -> epsilon
	   F -> (E)
	   F -> int
	*/

	TestG9 = parsergen.NewLL1Grammar(
		map[byte][]string{
			'P': {"E"},
			'E': {"TX"},
			'X': {"+TX", parsergen.Epsilon},
			'T': {"FU"},
			'U': {"*FU", parsergen.Epsilon},
			'F': {"(E)", parsergen.ValidLowercaseChar, parsergen.ValidInt},
		},
		append(append([]byte{'+', '*', '(', ')'}, parsergen.ValidLowercaseChar...), parsergen.ValidInt...), //terminals
		[]byte{'P', 'E', 'X', 'T', 'U', 'F'}, //nonterminals, in desired order
	)
)
