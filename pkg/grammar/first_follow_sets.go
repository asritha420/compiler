package grammar

import (
	"fmt"
)

/*
Given grammar symbol X

FIRST(X) is the set of all terminal symbols that can appear at the beginning of any string derived from X

For a terminal X:

	FIRST(X) = {X}

For a non terminal X with production rule X -> Y_1Y_2...Y_n

	For Y_i in Y_1Y_2...Y_n
		Add all terminal symbols in FIRST(Y_i) to FIRST(X) except epsilon
		If epsilon is in FIRST(Y_i), continue to Y_i+1
	If all Y_1, Y_2, ... Y_n can derive epsilon, add epsilon to FIRST(X)
*/

var (
	firstSets         = make(map[string][]rune) // the string can be any nonTerminal or terminal, not for ranges? -> should be symbol instead
	epsilonValue rune = ' '
)

func (g *Grammar) generateFirstSets() error {
	for _, rule := range g.Rules {
		firstSet, err := g.generateFirstSetFor(rule)
		if err != nil {
			return err
		}
		if existingFirstSet, ok := firstSets[rule.nonTerminal]; ok {
			firstSets[rule.nonTerminal] = append(existingFirstSet, firstSet...)
		} else {
			firstSets[rule.nonTerminal] = firstSet
		}
	}

	for k, v := range firstSets {
		fmt.Printf("FIRST SET for %s: \n", k)
		for _, r := range v {
			if r == epsilonValue {
				fmt.Printf(" EPSILON ")
			} else {
				fmt.Printf("%c", r)
			}
		}
		fmt.Println("")
	}

	return nil
}

func (g *Grammar) generateFirstSetFor(rule *Rule) ([]rune, error) {
	var ruleFirstSet []rune
	for _, prod := range rule.productions { // for each production within the Rule's productions list
	outerLoop:
		for i, s := range prod { // for each symbol within the current production
			switch s.symbolType {
			// if the symbol is a terminal
			case terminal, epsilon:
				terminalValue := s.validLiterals[0]
				firstLetter := []rune(terminalValue)[0]
				firstSets[terminalValue] = []rune{firstLetter}
				ruleFirstSet = append(ruleFirstSet, firstLetter) //should append
				break outerLoop
			case terminalLowercaseRange, terminalUppercaseRange, terminalNumberRange:
				var rangeValue []rune
				for _, validLiteral := range s.validLiterals {
					rangeValue = append(rangeValue, []rune(validLiteral)[0])
				}
				firstSets["RANGE"] = rangeValue //TODO: this is scuffed, also should be memoized
				ruleFirstSet = append(ruleFirstSet, rangeValue...)
				break outerLoop
			case nonTerminal:
				nT := s.validLiterals[0]
				var nTRule *Rule
				var nonTerminalFirstSet []rune

				// find corresponding rule for the nonTerminal
				for _, gRule := range g.Rules {
					if gRule.nonTerminal == nT {
						nTRule = gRule
					}
				}

				if nTRule == nil {
					return nil, fmt.Errorf("invalid non terminal: %s", nT)
				}

				// recursively generate the firstSet for that nonTerminal
				firstSet, err := g.generateFirstSetFor(nTRule)
				if err != nil {
					return nil, err
				}
				nonTerminalFirstSet = append(nonTerminalFirstSet, firstSet...)
				firstSets[nT] = nonTerminalFirstSet

				containsEpsilon := false

				// append all symbols in the nonTerminalFirstSet (FIRST(Y_i)) except epsilon
				for _, firstSetSymbol := range nonTerminalFirstSet {
					if firstSetSymbol == epsilonValue {
						containsEpsilon = true
					} else {
						ruleFirstSet = append(ruleFirstSet, firstSetSymbol)
					}
				}

				// continue to next symbol if it contains epsilon
				if !containsEpsilon {
					break outerLoop
				}

				// if all Y_1, Y_2, ... Y_n can derive epsilon, add epsilon to FIRST(X)
				if i == len(prod)-1 && containsEpsilon {
					ruleFirstSet = append(ruleFirstSet, epsilonValue)
					fmt.Printf("APPENDED EPSILON TO %v", rule.nonTerminal)
				}
			}
		}
	}
	return ruleFirstSet, nil
}

func (g *Grammar) generateFollowSets() {
	//TODO
}

// TODO: rename the symbol and production type to gSymbol, gProduction? so I don't hve to use letters in all for loops? or make sure all for loops have letters for consistency -> actually just use letters, just make sure that the same letters everywhere exist
//TODO: I don't like how the terminal and nonTerminal type validLiterals has a range, actually should be fine, its just that nonTerminal shouldn't be a range def ...//TODO: I don't like how validLiterals is a type []string, when it should be time []rune for the ranges, so split into two different things?
