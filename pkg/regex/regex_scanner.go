package regex

type RegexTokenType int

const (
	char RegexTokenType = iota

	// Alternation
	pipe

	// Quantifiers
	star     // Kleene star
	plus     // One or more
	question // Optional

	// Group
	openParenthesis
	closeParenthesis

	// Range
	openBracket
	closeBracket
	not
	dash
)

type RegexToken struct {
	RegexTokenType
	rune
}

func RegexScan(regexString string) []*RegexToken {
	tokens := make([]*RegexToken, 0)

	for _, r := range regexString {
		currToken := &RegexToken{}
		switch r {
		case '|':
			currToken.RegexTokenType = pipe
		case '*':
			currToken.RegexTokenType = star
		case '?':
			currToken.RegexTokenType = question
		case '(':
			currToken.RegexTokenType = openParenthesis
		case ')':
			currToken.RegexTokenType = closeParenthesis
		case '[':
			currToken.RegexTokenType = openBracket
		case ']':
			currToken.RegexTokenType = closeBracket
		case '^':
			currToken.RegexTokenType = not
		case '-':
			currToken.RegexTokenType = dash
		default:
			currToken.RegexTokenType = char
		}
		currToken.rune = r
		tokens = append(tokens, currToken)
	}

	return tokens
}
