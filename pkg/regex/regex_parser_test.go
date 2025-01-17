package regex

import (
	"testing"
)

/*
TO ADD in regex grammar:
- wildcard (.)
https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Regular_expressions/Cheatsheet
*/

// should test error cases as well
// TODO: use pointers for some of these nodes?
var (
	regex1            Regex = "a*"
	expectedRegex1AST       = KleeneStarNode{
		left: CharacterNode('a'),
	}

	regex2 Regex = "[abcd]"
	// TODO: need to fix so it rly returns this
	_                 = CharacterClassNode([]rune{'a', 'b', 'c', 'd'})
	expectedRegex2AST = AlternationNode{
		left: AlternationNode{
			left: AlternationNode{
				left:  CharacterNode('a'),
				right: CharacterNode('b'),
			},
			right: CharacterNode('c'),
		},
		right: CharacterNode('d'),
	}

	notRegex2           Regex = "[^abcd]"
	notRegex2Characters       = SubtractSlice2FromSlice1(anyChar, []rune{'a', 'b', 'c', 'd'})
	expectedNotRegex2         = CharacterClassNode(notRegex2Characters)

	// also do [a-zA-Z]? -> this doesnt make sense nvm
	// TODO: test epsilons
	regex3 Regex = "[a-zA-Z]+"
	// TODO: should equal this after I finish the refactor w/ alt/Character class
	_ = ConcatenationNode{
		left: CharacterClassNode(append(lowercaseLetters(), uppercaseLetters()...)),
		right: KleeneStarNode{
			left: CharacterClassNode(append(lowercaseLetters(), uppercaseLetters()...)),
		},
	}
	expectedRegex3AST = ConcatenationNode{
		left: AlternationNode{
			left:  CharacterClassNode(lowercaseLetters()),
			right: CharacterClassNode(uppercaseLetters()),
		},
		right: KleeneStarNode{
			left: AlternationNode{
				left:  CharacterClassNode(lowercaseLetters()),
				right: CharacterClassNode(uppercaseLetters()),
			},
		},
	}

	_      Regex = "[^a-zA-Z]?"
	regex6 Regex = "[e-]" // this is valid
)

func TestRegexParse(t *testing.T) {
	//// regex1
	//regex1AST, err := regex1.GetAST()
	//if err != nil {
	//	t.Errorf(err.Error())
	//}
	//
	//if regex1AST != expectedRegex1AST {
	//	t.Errorf("invalid")
	//}
	//
	//// regex2
	//regex2AST, err := regex2.GetAST()
	//if err != nil {
	//	t.Errorf(err.Error())
	//}
	//
	//if !reflect.DeepEqual(regex2AST, expectedRegex2AST) {
	//	t.Errorf("invalid")
	//}
	//
	//// notRegex2
	//notRegex2AST, err := notRegex2.GetAST()
	//if err != nil {
	//	t.Errorf(err.Error())
	//}
	//
	//if !reflect.DeepEqual(notRegex2AST, expectedNotRegex2) {
	//	t.Errorf("invalid")
	//}

	//// regex 3
	//regex3AST, err := regex3.GetAST()
	//if err != nil {
	//	t.Errorf(err.Error())
	//}
	//
	//if !reflect.DeepEqual(regex3AST, expectedRegex3AST) {
	//	t.Errorf("invalid")
	//}
}
