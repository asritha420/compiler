package regex

import (
	"reflect"
	"testing"
)

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

	regex3            Regex = "[a-zA-Z]?"
	expectedRegex3AST       = ConcatenationNode{
		left: CharacterClassNode(append(lowercaseLetters(), uppercaseLetters()...)),
		right: KleeneStarNode{
			left: CharacterClassNode(append(lowercaseLetters(), uppercaseLetters()...)),
		},
	}

	// TODO: is ^ outside of valid regex?
	_      Regex = "^[a-zA-Z]?"
	regex6 Regex = "[e-]"
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

	//// regex2
	//regex2AST, err := regex2.GetAST()
	//if err != nil {
	//	t.Errorf(err.Error())
	//}
	//
	//if !reflect.DeepEqual(regex2AST, expectedRegex2AST) {
	//	t.Errorf("invalid")
	//}

	// notRegex2
	// TODO: rename
	regex2AST, err := notRegex2.GetAST()
	if err != nil {
		t.Errorf(err.Error())
	}

	if !reflect.DeepEqual(regex2AST, expectedNotRegex2) {
		t.Errorf("invalid")
	}

	// regex 3
	//regex3AST, err := regex3.GetAST()
	//if err != nil {
	//	t.Errorf(err.Error()) // TODO: describe error
	//}
}
