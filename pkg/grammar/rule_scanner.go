package grammar

import (
	"fmt"
	"slices"
)

type production = []*symbol

type RuleScanner struct {
	curr                              int
	rule                              []rune
	validNonTerminals, validTerminals []string
}

func (rs *RuleScanner) Scan() ([]production, error) {
	productions := make([]production, 0)
	productions, err := rs.scanProduction(productions)
	if err != nil {
		return nil, fmt.Errorf("rs.Scan(): %v", err)
	}
	return productions, nil
}

func (rs *RuleScanner) scanProduction(productions []production) ([]production, error) {
	currProduction := make(production, 0)

	for rs.curr < len(rs.rule) {
		switch rs.rule[rs.curr] {
		case ' ':
			rs.curr++
		case '|':
			rs.curr++
			productions = append(productions, currProduction)
			return rs.scanProduction(productions)
		case '"':
			rs.curr++
			terminal, err := rs.consumeTerminal()
			if err != nil {
				return nil, err
			}
			currProduction = append(currProduction, newTerminalSymbol(terminal))
		case '[':
			rangeType, err := rs.consumeTerminalRange()
			if err != nil {
				return nil, err
			}
			rs.curr += 5 // consume the range
			currProduction = append(currProduction, newTerminalRangeSymbol(rangeType))
		default:
			nonTerminal, err := rs.consumeNonTerminal()
			if err != nil {
				return nil, err
			}
			currProduction = append(currProduction, newNonTerminalSymbol(nonTerminal))
		}
	}

	productions = append(productions, currProduction)
	return productions, nil
}

func (rs *RuleScanner) consumeNonTerminal() (string, error) {
	closingSpaceIndex := rs.consumeSymbolUntil(' ')

	var nT string
	if closingSpaceIndex == rs.curr { // the nonTerminal is the last symbol in the rule
		nT = string(rs.rule[rs.curr:len(rs.rule)])
	} else {
		nT = string(rs.rule[rs.curr:closingSpaceIndex])
	}

	if !rs.isValidNonTerminal(nT) {
		return "", fmt.Errorf("rs.consumeNonTerminal(): '%s' is an invalid nonTerminal", nT)
	}

	rs.curr += len(nT) // consume the non-terminal
	return nT, nil
}

func (rs *RuleScanner) consumeTerminal() (string, error) {
	closingQuoteIndex := rs.consumeSymbolUntil('"')

	t := string(rs.rule[rs.curr:closingQuoteIndex])

	if closingQuoteIndex == rs.curr {
		return "", fmt.Errorf("rs.consumeTerminal(): there is no closing quote for the terminal '%s'", t)
	}

	if !rs.isValidTerminal(t) && t != " " { //TODO: have a better way for handling epsilons?
		return "", fmt.Errorf("rs.consumeTerminal(): '%s' is an invalid terminal", t)
	}

	rs.curr = closingQuoteIndex + 1 // consume the terminal and closing quote
	return t, nil
}

// TODO: change to use a built in slices func to get first index of the end rune instead??
func (rs *RuleScanner) consumeSymbolUntil(end rune) int {
	endIndex := rs.curr

	for i := rs.curr; i < len(rs.rule); i++ {
		if rs.rule[i] == end {
			endIndex = i
			break
		}
	}

	return endIndex
}

func (rs *RuleScanner) consumeTerminalRange() (symbolType, error) {
	if !(len(rs.rule) >= rs.curr+4) {
		return -1, fmt.Errorf("rs.consumeTerminalRange(): rule '%s' ends midway through the terminal range, does not finish it", string(rs.rule))
	}

	terminalRange := string(rs.rule[rs.curr:(rs.curr + 5)])
	switch terminalRange {
	case "[a-z]":
		return terminalLowercaseRange, nil
	case "[A-Z]":
		return terminalUppercaseRange, nil
	case "[0-9]":
		return terminalNumberRange, nil
	default:
		return -1, fmt.Errorf("rs.consumeTerminalRange(): rule 	'%s' contains an invalid range definition: '%s'", string(rs.rule), terminalRange)
	}
}

func (rs *RuleScanner) isValidTerminal(candidate string) bool {
	return slices.Contains(rs.validTerminals, candidate)
}

func (rs *RuleScanner) isValidNonTerminal(candidate string) bool {
	return slices.Contains(rs.validNonTerminals, candidate)
}
