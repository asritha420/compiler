package regex

import "testing"

var (
	regex1String = "a*"
	regex1Tokens = []RegexToken {
		{
			RegexTokenType: char,
			rune: 'a'
		},
	}
	regexString2 = "[aaaa]"
	regexString3 = "^[a-zA-Z]?"
	// should test error cases as well
)


func TestRegexParse(t *testing.T) {
}
