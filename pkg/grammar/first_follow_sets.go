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

var (
	firstSets map[string][]rune       // the string can be any nonTerminal or terminal, not for ranges? -> should be symbol instead
	rulesMap  map[string][]production // need to generate this?
)

func (g *Grammar) generateFirstSets() error {
	for _, rule := range g.Rules {
		rulesMap[rule.nonTerminal] = rule.productions
	}

	for _, rule := range g.Rules {
		for _, p := range rule.productions { //should treat each individual rule prod as a rule?
			_, _ = g.generateFirstSetFOR(&p)
		}
	}

	for _, rule := range g.Rules {
		fmt.Printf("RULE for %s \n", rule.nonTerminal)
		fmt.Printf("%v", rule.FirstSet)
	}

	return nil
}

func (g *Grammar) generateFirstSetFOR(p *production) ([]rune, error) {
	var OGNTFirstset []rune
	for i, s := range *p {
		switch s.symbolType {
		case terminal:
			terminalValue := s.validLiterals[0]
			firstLetter := []rune(terminalValue)[0]
			firstSets[terminalValue] = []rune{firstLetter}
		case terminalLowercaseRange, terminalUppercaseRange, terminalNumberRange:
			var test []rune
			for _, validLiteral := range s.validLiterals {
				test = append(test, []rune(validLiteral)[0])
			}
		case nonTerminal:
			productions := rulesMap[s.validLiterals[0]]

			var nonTerminalFirstSet []rune // each first set for a rule/non-terminal will contain first sets from BOTH productions
			for _, p := range productions {
				firstSet, err := g.generateFirstSetFOR(&p)
				if err != nil {
					return nil, err
				}
				nonTerminalFirstSet = append(nonTerminalFirstSet, firstSet...)
			}

			containsEpsilon := false
			// TODO: def a better way to do this
			for _, lol := range nonTerminalFirstSet {
				if lol == ' ' { // append everything but the epsilon to the OGNTFIrstSet
					containsEpsilon = true // so go to the next one
				} else {
					OGNTFirstset = append(OGNTFirstset, lol)
				}
			}

			if !containsEpsilon {
				break
			}

			if i == len(*p)-1 && containsEpsilon {
				OGNTFirstset = append(OGNTFirstset, ' ')
			}
			// TODO: set it in the map so dont have to recalculate recursively? (ie. memoization) for the NTs
		}
	}
	return OGNTFirstset, nil
}

// recursively call itself? TODO: remainingSymbols should be pointer
func (g *Grammar) generateFirstSetFor(nT string) ([]rune, error) { // should be symbol instead of nT
	for _, p := range rulesMap[nT] {
		allEpsilon := true
		for _, s := range p {
			switch s.symbolType {
			case terminal:
				terminalValue := s.validLiterals[0]
				firstLetter := []rune(terminalValue)[0]
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

// TODO: rename the symbol and production type to gSymbol, gProduction? so I don't hve to use letters in all for loops? or make sure all for loops have letters for consistency -> actually just use letters, just make sure that the same letters everywhere exist

func (g *Grammar) generateFollowSets() {
	//TODO
}

//TODO: I don't like how the terminal and nonTerminal type validLiterals has a range, actually should be fine, its just that nonTerminal shouldn't be a range def ...//TODO: I don't like how validLiterals is a type []string, when it should be time []rune for the ranges, so split into two different things?
