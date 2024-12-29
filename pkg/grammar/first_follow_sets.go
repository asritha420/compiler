package grammar

import "fmt"

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
func (g *Grammar) generateFirstSets(s symbol) error {
	for _, rule := range g.Rules {
		rule.FirstSet = make([]rune, 0)
		for _, production := range rule.productions {
			for _, symbol := range production {
				switch symbol.symbolType {
				case terminal:
					rule.FirstSet = append(rule.FirstSet, ([]rune(symbol.validLiterals[0]))[0])
				case terminalLowercaseRange, terminalUppercaseRange, terminalNumberRange:
					for _, validLiteral := range symbol.validLiterals {
						rule.FirstSet = append(rule.FirstSet, []rune(validLiteral)[0])
					}
				case nonTerminal:
				/*
					get firstSet for the nonTerminal recursively
					if that firstSet contains epsilon, proceed to next iteration
				*/

				case epsilon:

				default:
					return fmt.Errorf("g.generateFirstSets(): symbol %v has invalid symbol type") // TODO: fix
				}
			}
		}
	}
	return nil
}

func (g *Grammar) generateFirstSetsHelper(s symbol) ([]rune, error) {
	return nil, nil
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
