package parsergen

import (
	"slices"
)

// TODO: this should be of type LL1LL1Grammar?
// TODO: include support for nonterminalprimes in the grammar rule names

/*
an LL(1) grammar:
L: Left-to-right scanning
L: leftmost derivation of the input string
1: at most one symbol lookahead used
*/
type LL1Grammar struct {
	Rules        map[byte][]string
	terminals    []byte
	nonterminals []byte
	FirstSet     map[byte][]byte
	FollowSet    map[byte][]byte
}

const (
	Epsilon   = ""
	eof       = '$'
	ValidInt  = "0123456789"
	ValidChar = "abcdefghijklmnopqrstuvwxyz"
)

func NewLL1Grammar(Rules map[byte][]string) *LL1Grammar {
	g := LL1Grammar{
		Rules:     Rules,
		FirstSet:  map[byte][]byte{},
		FollowSet: map[byte][]byte{},
	}

	var nonterminals []byte
	var terminals []byte

	for nt, _ := range g.Rules {
		nonterminals = append(nonterminals, nt)
	}

	for _, rules := range g.Rules {
		for _, r := range rules {
			if r == ValidInt {
				terminals = append(terminals, ValidInt...)
			} else if r == ValidChar {
				terminals = append(terminals, ValidChar...)
			} else {
				for _, b := range []byte(r) {
					if !slices.Contains(nonterminals, b) {
						terminals = append(terminals, b)
					}
				}
			}
		}
	}

	g.terminals = terminals
	g.nonterminals = nonterminals

	for _, startNt := range g.nonterminals {
		g.generateFirstSet(startNt, startNt)
		//		g.generateFollowSet(startNt, startNt)
	}

	return &g
}

// generateFirstSet() TODO: fill in description
// if its an Epsilon, append it
// if the first character of the Rules[rName] is a terminal, then append it.
// if the first character of the Rules[rName] is a nonTerminal, recursively call follow until it reaches a terminal
func (g *LL1Grammar) generateFirstSet(ogByte byte, nt byte) {
	for _, r := range g.Rules[nt] {
		if r == Epsilon {
			g.FirstSet[nt] = append(g.FirstSet[nt], 0x00)
		} else if r == ValidChar {
			g.FirstSet[nt] = append(g.FirstSet[nt], ValidChar...)
			g.FirstSet[ogByte] = append(g.FirstSet[ogByte], ValidChar...)
		} else if r == ValidInt {
			g.FirstSet[nt] = append(g.FirstSet[nt], ValidInt...)
			g.FirstSet[ogByte] = append(g.FirstSet[ogByte], ValidInt...)
		} else if slices.Contains(g.terminals, r[0]) {
			g.FirstSet[nt] = append(g.FirstSet[nt], r[0])
			g.FirstSet[ogByte] = append(g.FirstSet[ogByte], r[0]) //only add if nt != ogByte
		} else if slices.Contains(g.nonterminals, r[0]) {
			g.generateFirstSet(ogByte, r[0])
		}
	}
}

//TODO: better naming than "ogByte"

// TODO: i dont have to loop through the nonTerminals, i think starting from startByte will fill everything in
// TODO: right now, the naming is in inconsistent, im using rName vs ntbyte
// for rule in Rules[rName]list
func (g *LL1Grammar) generateFollowSet(ogByte byte, nt byte) {
	/*
		if nt == S, append $ //TODO: is $ called eof?

		else:
			recursively go through the last character in the r.
				if it is a terminal. add it to the list + og byte
				else if it is a nonterminal (beta):
					recursively call generateFollowSet(ogByte, nt byte) //this will generate FOLLOW(A)
					if there is a nonterminal in the second to last position (B):
						if firstSets[beta] contains Epsilon:
							followSets[B].append all elements within firstSets[beta] (except for Epsilon)



	*/
	if nt == g.nonterminals[0] {
		g.FollowSet[nt] = append(g.FollowSet[nt], eof)
	} else {
		for _, r := range g.Rules[nt] {
			if slices.Contains(g.terminals, r[len(r)-1]) { //if last character is a terminal
				g.FollowSet[nt] = append(g.FollowSet[nt], r[len(r)-1])
				g.FollowSet[ogByte] = append(g.FollowSet[ogByte], r[len(r)-1])
			} else if slices.Contains(g.nonterminals, r[len(r)-1]) {
				if slices.Contains(g.Rules[r[len(r)-1]], Epsilon) {
					//add FIRST(second to last) to FOLLOW B
					//what does expecting Epsilon mean?
				}
				//do generateFollowSet(ogByte, second to last terminal)
				//then add every element followSet[ogByte] to followSet[second to last terminal]
			}
		}
	}
}
