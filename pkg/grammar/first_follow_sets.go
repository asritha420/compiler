package grammar

import (
	"fmt"
	"slices"
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

// TODO: currently doesnt do EXCEPT epsilon on line 20
var (
	firstSets map[string][]rune //the string can be any nonTerminal or terminal, not for ranges?
	rulesMap  map[string][]production
)

func (g *Grammar) generateFirstSets() error {
	for _, rule := range g.Rules {
		rulesMap[rule.nonTerminal] = rule.productions
	}

	for _, rule := range g.Rules {
		for _, p := range rule.productions { //should treat each individual rule prod as a rule?
			for _, s := range p {
				switch s.symbolType {
				case terminal:
					terminalValue := s.validLiterals[0]
					firstLetter := []rune(terminalValue)[0] //TODO: this is messy lol
					firstSets[terminalValue] = []rune{firstLetter}
				case terminalLowercaseRange, terminalUppercaseRange, terminalNumberRange:
					var test []rune
					for _, validLiteral := range s.validLiterals {
						test = append(test, []rune(validLiteral)[0])
					}
				case nonTerminal:
					g.generateFirstSetFor(s.validLiterals[0]) // not doing anything with the returned
				}
			}
		}
	}

	for _, rule := range g.Rules {
		fmt.Printf("RULE for %s \n", rule.nonTerminal)
		fmt.Printf("%v", rule.FirstSet)
	}

	return nil
}

// recursively call itself? TODO: remainingSymbols should be pointer
func (g *Grammar) generateFirstSetFor(nT string) ([]rune, error) {
	for _, p := range rulesMap[nT] {
		allEpsilon := true
		for _, s := range p {
			switch s.symbolType {
			case terminal:
				terminalValue := s.validLiterals[0]
				firstLetter := []rune(terminalValue)[0] //TODO: this is messy lol
				firstSets[terminalValue] = []rune{firstLetter}
				return firstSets[terminalValue], nil
			case terminalLowercaseRange, terminalUppercaseRange, terminalNumberRange:
				var test []rune
				for _, validLiteral := range s.validLiterals {
					test = append(test, []rune(validLiteral)[0])
				}
				return test, nil
			case nonTerminal:
				/*
					get firstSet for the nonTerminal recursively
					if that firstSet contains epsilon, proceed to next iteration
				*/
				firstSet, err := g.generateFirstSetFor(s.validLiterals[0])
				if err != nil {
					return nil, err
				}
				if value, ok := firstSets[nT]; ok {
					firstSets[nT] = append(value, firstSet...)
				} else {
					firstSets[nT] = firstSet
				}
				/*
					- for the rest of the symbols in the curr production left, keep going until none of them have epsilon, add as you keep going
					- if all of them have epsilon, add epsilon to firstX //TODO: should have a getNewEpsilon()? instead of using getNewTerminal?
				*/
				if !slices.Contains(firstSet, ' ') { //one of them does not contain an epsilon
					allEpsilon = false
					break // stop the for loop
				}
			case epsilon:
				return []rune{' '}, nil
			default:
				return nil, fmt.Errorf("g.generateFirstSets(): symbol has invalid symbol type") // TODO: fix
			}
		}
		if allEpsilon {
			if value, ok := firstSets[nT]; ok {
				firstSets[nT] = append(value, ' ')
			} else {
				firstSets[nT] = []rune{' '}
			}
		}
	}
	return firstSets[nT], nil
}

// TODO: rename the symbol and production type to gSymbol, gProduction? so I odnt hve to use letters in all for loops? or make sure all for loops have letters for consistency

func (g *Grammar) generateFollowSets() {
	//TODO
}

//TODO: I don't like how the terminal and nonTerminal type validLiterals has a range, actually should be fine, its just that nonTerminal shouldn't be a range def ...//TODO: I don't like how validLiterals is a type []string, when it should be time []rune for the ranges, so split into two different things?
