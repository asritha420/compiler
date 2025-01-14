package regex

import "testing"

var (
	regex1String = "a*"
	regex1Tokens = []RegexToken{
		{
			RegexTokenType: char,
		},
	}
	regex1AST = KleeneStar{
		left: Const{'a'},
	}
	regexString2 = "[abcd]"
	regex2AST    = Alternation{
		left: Const{'a'},
	}
	regexString3 = "^[a-zA-Z]?"
	// should test error cases as well
	omsFunkyRegex    = "r(e|t*?t)[^hel-p]*"
	omsFunkyRegexAST = Concatenation{
		left: Const{'r'},
		right: Alternation{
			left: Const{'e'},
			right: KleeneStar{
				left: Const{'t'},
			},
		},
	}
)

/*
[asdf]
	alt
	alt f
  alt d
alt s
alt a
*/

func TestRegexParse(t *testing.T) {
}
