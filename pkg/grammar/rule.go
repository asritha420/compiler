package grammar

type Rule struct {
	nonTerminal         string
	productions         []production
	FirstSet, FollowSet []rune
}



