package parsergen

/*
an LL(1) grammar:
L: Left-to-right scanning
L: leftmost derivation of the input string
1: at most one symbol lookahead used
*/
//
//
//type LL1Grammar struct {
//	*grammar
//}
//
//func NewLL1Grammar(rules []Rule) *LL1Grammar {
//	ll1g := &LL1Grammar{
//		grammar: newGrammar(rules),
//	}
//
//	return ll1g
//}
//
//////TODO: clean up both the first and follow functions as they are absolute messes, also should they be methods on an interface Grammar? or should I use composition, look ahead in the book or ask Chatgpt to determine if there are different ways of calculating first and follow sets for different types of grammars
////
////// rules: map[byte][]string
////func (g *LL1Grammar) generateFollowSet(rName byte) {
////	//if the startRName byte is = to the start symbol, add $ to the follow set
////	startSymbol := g.nonTerminals[0]
////	if rName == startSymbol {
////		g.FollowSet[rName] = append(g.FollowSet[rName], endOfInput)
////	}
////
////	//TODO: move it so that before it loops through each rule in grammar nad calls this generateFollowSet function, it FIRST calls the helper addTerminalsHelper just ONCE
////
////	for _, p := range g.Rules[rName] {
////		if p == Epsilon {
////			continue
////		}
////
////		//TODO: this isnt checking if the p matches any of the two valid forms correctly
////		//TODO: fix so not using capital B, as per Go naming conventions
////		if len(p) > 1 {
////			for i := 1; i < len(p); i++ {
////				if slices.Contains(g.nonTerminals, p[i]) { // find the first occuring nonterminal B
////					B := p[i]
////					if i != len(p)-1 { //there is a beta
////						beta := p[i+1 : len(p)]
////						betaFirstSet := g.getFirstSetHelper(beta)
////						g.FollowSet[B] = append(g.FollowSet[B], betaFirstSet...)
////						if slices.Contains(betaFirstSet, 0x00) { //if betaFirstSet contains epsilon
////							g.addFollowTerminalsHelper(rName)
////							g.FollowSet[B] = append(g.FollowSet[B], g.FollowSet[rName]...)
////							break //break out of this for loop
////						}
////					}
////					g.addFollowTerminalsHelper(rName)
////					g.FollowSet[B] = append(g.FollowSet[B], g.FollowSet[rName]...)
////					break
////				}
////			}
////		}
////
////		/*
////			given a string p
////
////			check if its aB
////			or aBbeta
////			where B is a nonterminal
////
////			is string is longer than 1 character
////			keep going in the string until it finds a valid nonTerminal B
////			let beta = the substring after the nonTerminal B if B is not the last character in the string
////		*/
////	}
////}
////
////// TODO: probably get rid of this or incorporate it into the actual getFirstSet() logic
////func (g *LL1Grammar) getFirstSetHelper(beta string) []byte {
////	if len(beta) == 1 {
////		symbol := byte(beta[0])
////		if slices.Contains(g.terminals, symbol) {
////			return []byte{symbol}
////		} else if slices.Contains(g.nonTerminals, symbol) {
////			return g.FirstSet[symbol]
////		}
////	}
////
////	//TODO: implement logic for when the beta is more than one character long, rn it just returns nil, I think it should include for the last byte of the string, and then return that
////	return []byte{}
////}
////
////// TODO: probably don't need this
////// TODO: change this to rName?
////func (g *LL1Grammar) addFollowTerminalsHelper(nt byte) {
////	var followingTerminals []byte
////
////	for _, productions := range g.Rules {
////		for _, p := range productions {
////			if strings.Contains(p, string(nt)) {
////				for i := 0; i < len(p)-1; i++ { //loop through all bytes in the production except the last byte
////					if p[i] == nt { //if the byte is a nonTerminal and the following byte is a nonTerminal
////						if slices.Contains(g.terminals, p[i+1]) {
////							followingTerminals = append(followingTerminals, p[i+1])
////						}
////					}
////				}
////			}
////		}
////	}
////
////	g.FollowSet[nt] = append(g.FollowSet[nt], followingTerminals...)
////}
////
