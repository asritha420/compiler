package grammar

/*
Given grammar symbol X

# FIRST(X) is the set of all terminal symbols that can appear at the beginning of any string derived from X

For a terminal X:

	FIRST(X) = {X}

For a production rule X -> Y_1Y_2...Y_n

	For Y_i in Y_1Y_2...Y_n
		Add all terminal symbols in FIRST(Y_i) to FIRST(X) except epsilon
		If epsilon is in FIRST(Y_i), continue to Y_i+1
	If all Y_1, Y_2, ... Y_n can derive epsilon, add epsilon to FIRST(X)
*/
func (g *Grammar) generateFirstSets() {

	//for i := len(g.Rules) - 1; i >= 0; i-- {
	//	r := g.Rules[i]
	//
	//}
}

func (g *Grammar) generateFollowSets() {
	//TODO
}

//
//// generateFirstSet() will recursively calculate the FIRST set for each rule in the Grammar;  FIRST(x) is the set of all terminals (including epsilon) that can appear at the beginning of any possible derivation of rule x.
//func (g *Grammar) generateFirstSet(startRName byte, rName byte) {
//	for _, r := range g.Rules[rName] {
//		//if the rule is an epsilon or a valid character or int, add it to FIRST(rName)
//		if r == Epsilon {
//			g.FirstSet[rName] = append(g.FirstSet[rName], 0x00)
//		} else if r == ValidChar {
//			g.FirstSet[rName] = append(g.FirstSet[rName], ValidChar...)
//			g.FirstSet[startRName] = append(g.FirstSet[startRName], ValidChar...)
//		} else if r == ValidInt {
//			g.FirstSet[rName] = append(g.FirstSet[rName], ValidInt...)
//			g.FirstSet[startRName] = append(g.FirstSet[startRName], ValidInt...)
//		} else if slices.Contains(g.terminals, r[0]) {
//			g.FirstSet[rName] = append(g.FirstSet[rName], r[0]) // if the first byte of Rules[rName] is a terminal, add it to FIRST(rName)
//			g.FirstSet[startRName] = append(g.FirstSet[startRName], r[0])
//		} else if slices.Contains(g.nonTerminals, r[0]) {
//			g.generateFirstSet(startRName, r[0]) // if the first byte of Rules[rName] is a nonTerminal, recursively call generateFirstSet() until it finds a terminal
//		}
//	}
//}
//func (g *Grammar) generateFirstSet(rule *Rule, nextRule *Rule) {
//
//	for _, rule := range g.Rules {
//		for _, production := range rule.Productions {
//			if production == Epsilon {
//				rule.firstSet = append(rule.firstSet, Epsilon)
//			} else if production == ValidInt || production == ValidChar {
//				rule.firstSet = append(rule.firstSet, production)
//			}
//		}
//	}
//}
