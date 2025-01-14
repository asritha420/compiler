package regex

import "testing"

// should test error cases as well
var (
	regex1    = "a*"
	regex1AST = KleeneStarNode{
		left: LiteralNode('a'),
	}

	regex2    = "[abcd]"
	regex2AST = CharacterClassNode([]rune{'a', 'b', 'c', 'd'})

	regex3    = "^[a-zA-Z]?"
	regex3AST = CharacterClassNode([]rune{})

	regex4 = "r(e|t*?t)[^hel-p]*"
)

func TestRegexParse(t *testing.T) {
}
