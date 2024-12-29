package grammar

type Rule struct {
	nonTerminal         string
	productions         [][]*symbol
	FirstSet, FollowSet []rune
}
