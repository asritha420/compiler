package grammar

var (
	regexGrammar = [][]rune{
		[]rune(`production = expression`),
		[]rune(`expression = term expressionPrime`),
		[]rune(`expressionPrime = "|" term expressionPrime | EPSILON`),
		[]rune(`term = factor termPrime`),
		[]rune(`termPrime = factor termPrime | EPSILON`),
		[]rune(`factor = group factorPrime`),
		[]rune(`factorPrime = "*" factorPrime | EPSILON`),
		[]rune(`group = "(" expression ")" | [a-z] | [A-Z] | [0-9]`),
	}
)
