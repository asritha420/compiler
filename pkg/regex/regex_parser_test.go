package regex

import "testing"

// should test error cases as well
// TODO: use pointers for some of these nodes?
var (
	regex1            Regex = "a*"
	expectedRegex1AST       = KleeneStarNode{
		left: CharacterNode('a'),
	}

	regex2            Regex = "[abcd]"
	expectedRegex2AST       = CharacterClassNode([]rune{'a', 'b', 'c', 'd'})

	regex3            Regex = "[a-zA-Z]?"
	expectedRegex3AST       = ConcatenationNode{
		left: CharacterClassNode(append(lowercaseLetters(), uppercaseLetters()...)),
		right: KleeneStarNode{
			left: CharacterClassNode(append(lowercaseLetters(), uppercaseLetters()...)),
		},
	}

	_      Regex = "^[a-zA-Z]?"
	regex6 Regex = "[e-]"
)

func TestRegexParse(t *testing.T) {

	regex1AST, err := regex1.GetAST()
	if err != nil {
		t.Errorf(err.Error())
	}

	if regex1AST != expectedRegex1AST {
		t.Errorf("invalid")
	}
}
