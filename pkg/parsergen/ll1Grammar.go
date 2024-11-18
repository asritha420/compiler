package parsergen

import (
	"fmt"
	"slices"
	"strings"
)

/*
an LL(1) grammar:
L: Left-to-right scanning
L: leftmost derivation of the input string
1: at most one symbol lookahead used
*/

type LL1Grammar struct {
	Rules        map[byte][]string //Rules can be an object
	terminals    []byte
	nonTerminals []byte
	FirstSet     map[byte][]byte
	FollowSet    map[byte][]byte
}

const (
	Epsilon    = ""
	endOfInput = '$'
	//	ValidInt           = "0123456789" TODO: uncomment
	//	ValidLowercaseChar = "abcdefghijklmnopqrstuvwxyz" //TODO: combine into one
	ValidCapitalChar   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	ValidInt           = "0"
	ValidLowercaseChar = "a"
)

// TODO: see below, make it so that you don't have to specify terminals and nonTerminals
func NewLL1Grammar(Rules map[byte][]string, terminals []byte, nonTerminals []byte) *LL1Grammar {
	g := LL1Grammar{
		Rules:        Rules,
		FirstSet:     map[byte][]byte{},
		FollowSet:    map[byte][]byte{},
		terminals:    terminals,
		nonTerminals: nonTerminals,
	}

	//var nonTerminals []byte
	//var terminals []byte
	//
	////TODO: look into maps.Keys(<mapName>)
	//for nt, _ := range g.Rules {
	//	nonTerminals = append(nonTerminals, nt)
	//}

	// fills in terminals and nonTerminals slices #TODO: this wont fill it in the same order every time, bc when loop through map is random order
	//for _, rules := range g.Rules {
	//	for _, r := range rules {
	//		if r == ValidInt {
	//			terminals = append(terminals, ValidInt...)
	//		} else if r == ValidLowercaseChar {
	//			terminals = append(terminals, ValidLowercaseChar...)
	//		} else {
	//			for _, b := range []byte(r) {
	//				if !slices.Contains(nonTerminals, b) {
	//					terminals = append(terminals, b)
	//				}
	//			}
	//		}
	//	}
	//}
	//
	//g.terminals = terminals
	//g.nonTerminals = nonTerminals

	for _, name := range g.nonTerminals {
		g.generateFirstSet(name, name)
		//g.generateFollowSet(name)
	}

	for _, name := range g.nonTerminals {
		fmt.Printf("NT: %v %c \n", name, name)
	}

	for _, name := range g.nonTerminals {
		g.generateFollowSet(name)
	}

	return &g
}

// generateFirstSet() will recursively calculate the FIRST set for each rule in the grammar;  FIRST(x) is the set of all terminals (including epsilon) that can appear at the beginning of any possible derivation of rule x.
func (g *LL1Grammar) generateFirstSet(startRName byte, rName byte) {
	for _, r := range g.Rules[rName] {
		//if the rule is an epsilon or a valid character or int, add it to FIRST(rName)
		if r == Epsilon {
			g.FirstSet[rName] = append(g.FirstSet[rName], 0x00)
		} else if r == ValidLowercaseChar {
			g.FirstSet[rName] = append(g.FirstSet[rName], ValidLowercaseChar...)
			g.FirstSet[startRName] = append(g.FirstSet[startRName], ValidLowercaseChar...)
		} else if r == ValidInt {
			g.FirstSet[rName] = append(g.FirstSet[rName], ValidInt...)
			g.FirstSet[startRName] = append(g.FirstSet[startRName], ValidInt...)
		} else if slices.Contains(g.terminals, r[0]) {
			g.FirstSet[rName] = append(g.FirstSet[rName], r[0]) // if the first byte of Rules[rName] is a terminal, add it to FIRST(rName)
			g.FirstSet[startRName] = append(g.FirstSet[startRName], r[0])
		} else if slices.Contains(g.nonTerminals, r[0]) {
			g.generateFirstSet(startRName, r[0]) // if the first byte of Rules[rName] is a nonTerminal, recursively call generateFirstSet() until it finds a terminal
		}
	}
}

//TODO: clean up both the first and follow functions as they are absolute messes, also should they be methods on an interface Grammar? or should I use composition, look ahead in the book or ask Chatgpt to determine if there are different ways of calculating first and follow sets for different types of grammars

// rules: map[byte][]string
func (g *LL1Grammar) generateFollowSet(rName byte) {
	//if the startRName byte is = to the start symbol, add $ to the follow set
	startSymbol := g.nonTerminals[0]
	if rName == startSymbol {
		g.FollowSet[rName] = append(g.FollowSet[rName], endOfInput)
	}

	//TODO: move it so that before it loops through each rule in grammar nad calls this generateFollowSet function, it FIRST calls the helper addTerminalsHelper just ONCE

	for _, p := range g.Rules[rName] {
		if p == Epsilon {
			continue
		}

		//TODO: this isnt checking if the p matches any of the two valid forms correctly
		//TODO: fix so not using capital B, as per Go naming conventions
		if len(p) > 1 {
			for i := 1; i < len(p); i++ {
				if slices.Contains(g.nonTerminals, p[i]) { // find the first occuring nonterminal B
					B := p[i]
					if i != len(p)-1 { //there is a beta
						beta := p[i+1 : len(p)]
						betaFirstSet := g.getFirstSetHelper(beta)
						g.FollowSet[B] = append(g.FollowSet[B], betaFirstSet...)
						if slices.Contains(betaFirstSet, 0x00) { //if betaFirstSet contains epsilon
							g.addFollowTerminalsHelper(rName)
							g.FollowSet[B] = append(g.FollowSet[B], g.FollowSet[rName]...)
							break //break out of this for loop
						}
					}
					g.addFollowTerminalsHelper(rName)
					g.FollowSet[B] = append(g.FollowSet[B], g.FollowSet[rName]...)
					break
				}
			}
		}

		/*
			given a string p

			check if its aB
			or aBbeta
			where B is a nonterminal

			is string is longer than 1 character
			keep going in the string until it finds a valid nonTerminal B
			let beta = the substring after the nonTerminal B if B is not the last character in the string
		*/
	}
}

// TODO: probably get rid of this or incorporate it into the actual getFirstSet() logic
func (g *LL1Grammar) getFirstSetHelper(beta string) []byte {
	if len(beta) == 1 {
		symbol := byte(beta[0])
		if slices.Contains(g.terminals, symbol) {
			return []byte{symbol}
		} else if slices.Contains(g.nonTerminals, symbol) {
			return g.FirstSet[symbol]
		}
	}

	//TODO: implement logic for when the beta is more than one character long, rn it just returns nil, I think it should include for the last byte of the string, and then return that
	return []byte{}
}

// TODO: probably don't need this
// TODO: change this to rName?
func (g *LL1Grammar) addFollowTerminalsHelper(nt byte) {
	var followingTerminals []byte

	for _, productions := range g.Rules {
		for _, p := range productions {
			if strings.Contains(p, string(nt)) {
				for i := 0; i < len(p)-1; i++ { //loop through all bytes in the production except the last byte
					if p[i] == nt { //if the byte is a nonTerminal and the following byte is a nonTerminal
						if slices.Contains(g.terminals, p[i+1]) {
							followingTerminals = append(followingTerminals, p[i+1])
						}
					}
				}
			}
		}
	}

	g.FollowSet[nt] = append(g.FollowSet[nt], followingTerminals...)
}
