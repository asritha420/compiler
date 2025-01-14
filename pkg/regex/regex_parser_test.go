package regex

import "testing"

var (
	regexString1 = "a*"
	regexString2 = "[aaaa]"
	regexString3 = "^[a-zA-Z]?"
	// should test error cases as well
)

func TestRegexParse(t *testing.T) {
}
