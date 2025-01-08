package main

// import (
// 	"log"
// 	"os"
// 	"reflect"
// )

// type token struct {
// 	nonTerm string
// 	literal []rune
// 	tree    []*token
// }

// type ebnfAST interface {
// 	match([]rune) *token
// 	/*
// 		returns true if the input could be the valid start of a token
// 	*/
// 	canMatch([]rune) bool
// }

// type ebnfNonTerm struct {
// 	nonTerm string
// 	rule    ebnfAST
// }

// func (e *ebnfNonTerm) match(s []rune) *token {
// 	if t := e.rule.match(s); t != nil {
// 		t.nonTerm = e.nonTerm
// 		return t
// 	}
// 	return nil
// }

// func (e *ebnfNonTerm) canMatch(s []rune) bool {
// 	return e.rule.canMatch(s)
// }

// func newNonTerm(nonTerm string, rule ebnfAST) *ebnfNonTerm {
// 	return &ebnfNonTerm{nonTerm: nonTerm, rule: rule}
// }

// type ebnfTerm struct {
// 	term []rune
// }

// func (e *ebnfTerm) match(s []rune) *token {
// 	if reflect.DeepEqual(e.term, s) {
// 		return &token{literal: s}
// 	}
// 	return nil
// }

// func (e *ebnfTerm) canMatch(s []rune) bool {
// 	if len(s) > len(e.term) {
// 		return false
// 	}

// 	for i, c := range s {
// 		if c != e.term[i] {
// 			return false
// 		}
// 	}

// 	return true
// }

// func newTerm(term string) *ebnfTerm {
// 	return &ebnfTerm{term: []rune(term)}
// }

// type ebnfToRange struct {
// 	low, high rune
// }

// func (e *ebnfToRange) match(s []rune) *token {
// 	if len(s) != 1 {
// 		return nil
// 	}

// 	if s[0] >= e.low && s[0] < e.high {
// 		return &token{literal: s}
// 	}
// 	return nil
// 	// not inclusive of high
// }

// func (e *ebnfToRange) canMatch(s []rune) bool {

// 	return len(s) == 0 || (len(s) == 1 && s[0] >= e.low && s[0] < e.high)
// }

// func newToRange(low, high rune) *ebnfToRange {
// 	return &ebnfToRange{low: low, high: high}
// }

// type ebnfCharRange struct {
// 	validChars map[rune]struct{}
// }

// func (e *ebnfCharRange) match(s []rune) *token {
// 	if len(s) != 1 {
// 		return nil
// 	}

// 	if _, ok := e.validChars[s[0]]; ok {
// 		return &token{literal: s}
// 	}

// 	return nil
// }

// func (e *ebnfCharRange) canMatch(s []rune) bool {
// 	if len(s) != 1 {
// 		return len(s) == 0
// 	}

// 	_, ok := e.validChars[s[0]]

// 	return ok
// }

// func newCharRange(validChars string) *ebnfCharRange {
// 	runeMap := make(map[rune]struct{})
// 	for _, c := range validChars {
// 		runeMap[c] = struct{}{}
// 	}
// 	return &ebnfCharRange{validChars: runeMap}
// }

// type ebnfConcat struct {
// 	children []ebnfAST
// }

// func (e *ebnfConcat) match(s []rune) *token {
// 	lastMatchIdx := 0
// 	currIdx := 0

// 	tokens := make([]*token, 0)

// 	for _, ebnf := range e.children {
// 		if lastMatchIdx == len(s) {
// 			// end of string (safety against oob reads)
// 			// check if it can match against empty string
// 			if ebnf.match(make([]rune, 0)) == nil {
// 				return nil
// 			}
// 			continue
// 		}

// 		// go until token can't match or end of string
// 		for ebnf.canMatch(s[lastMatchIdx:currIdx]) {
// 			currIdx++
// 			if currIdx > len(s) {
// 				break
// 			}
// 		}

// 		// back off one to possibly get valid token
// 		// Note it could be incomplete!
// 		currIdx--

// 		// create token
// 		token := ebnf.match(s[lastMatchIdx:currIdx])
// 		if token == nil {
// 			return nil
// 		}

// 		tokens = append(tokens, token)
// 		lastMatchIdx = currIdx
// 	}

// 	if lastMatchIdx != len(s) {
// 		return nil
// 	}

// 	return &token{literal: s, tree: tokens}
// }

// func (e *ebnfConcat) canMatch(s []rune) bool {
// 	lastMatchIdx := 0
// 	currIdx := 0

// 	for _, ebnf := range e.children {
// 		// go until token can't match or end of string
// 		for ebnf.canMatch(s[lastMatchIdx:currIdx]) {
// 			currIdx++
// 			if currIdx > len(s) {
// 				return true
// 			}
// 		}

// 		// back off one to possibly get valid token
// 		// Note it could be incomplete!
// 		currIdx--

// 		// check if it matches
// 		if ebnf.match(s[lastMatchIdx:currIdx]) == nil {
// 			return false
// 		}

// 		lastMatchIdx = currIdx
// 	}

// 	return lastMatchIdx == len(s)
// }

// func newConcat(children ...ebnfAST) *ebnfConcat {
// 	return &ebnfConcat{children: children}
// }

// type ebnfAlt struct {
// 	children []ebnfAST
// }

// func (e *ebnfAlt) match(s []rune) *token {
// 	for _, child := range e.children {
// 		if t := child.match(s); t != nil {
// 			return t
// 		}
// 	}

// 	return nil
// }

// func (e *ebnfAlt) canMatch(s []rune) bool {
// 	for _, child := range e.children {
// 		if child.canMatch(s) {
// 			return true
// 		}
// 	}

// 	return false
// }

// func newAlt(children ...ebnfAST) *ebnfAlt {
// 	return &ebnfAlt{children: children}
// }

// type ebnfException struct {
// 	rule, except ebnfAST
// }

// func (e *ebnfException) match(s []rune) *token {
// 	t := e.rule.match(s)
// 	if t == nil {
// 		return nil
// 	}

// 	if e.except.match(s) != nil {
// 		return nil
// 	}

// 	return t
// }

// func (e *ebnfException) canMatch(s []rune) bool {
// 	return e.rule.canMatch(s) && e.except.match(s) == nil
// }

// func newException(rule, except ebnfAST) *ebnfException {
// 	return &ebnfException{rule: rule, except: except}
// }

// type ebnfNoneOrOne struct {
// 	child ebnfAST
// }

// func (e *ebnfNoneOrOne) match(s []rune) *token {
// 	if len(s) == 0 {
// 		return &token{literal: s}
// 	}

// 	return e.child.match(s)
// }

// func (e *ebnfNoneOrOne) canMatch(s []rune) bool {
// 	return e.child.canMatch(s)
// }

// func newNoneOrOne(child ebnfAST) *ebnfNoneOrOne {
// 	return &ebnfNoneOrOne{child: child}
// }

// type ebnfNoneOrMore struct {
// 	child ebnfAST
// }

// func (e *ebnfNoneOrMore) match(s []rune) *token {
// 	if len(s) == 0 {
// 		return &token{literal: s}
// 	}

// 	lastMatchIdx := 0
// 	currIdx := 0

// 	tokens := make([]*token, 0)

// 	for {
// 		// go until token can't match or end of string
// 		for e.child.canMatch(s[lastMatchIdx:currIdx]) {
// 			currIdx++
// 			if currIdx > len(s) {
// 				break
// 			}
// 		}

// 		// back off one to possibly get valid token
// 		// Note it could be incomplete!
// 		currIdx--

// 		// create token
// 		token := e.child.match(s[lastMatchIdx:currIdx])
// 		if token == nil {
// 			return nil
// 		}
// 		tokens = append(tokens, token)

// 		if currIdx == len(s) {
// 			break
// 		}
// 		lastMatchIdx = currIdx
// 	}

// 	return &token{literal: s, tree: tokens}
// }

// func (e *ebnfNoneOrMore) canMatch(s []rune) bool {
// 	lastMatchIdx := 0
// 	currIdx := 0

// 	for {
// 		// go until token can't match or end of string
// 		for e.child.canMatch(s[lastMatchIdx:currIdx]) {
// 			currIdx++
// 			if currIdx > len(s) {
// 				return true
// 			}
// 		}

// 		// back off one to possibly get valid token
// 		// Note it could be incomplete!
// 		currIdx--

// 		if e.child.match(s[lastMatchIdx:currIdx]) == nil {
// 			return false
// 		}

// 		if currIdx == len(s) {
// 			break
// 		}
// 		lastMatchIdx = currIdx
// 	}

// 	return true
// }

// func newNoneOrMore(child ebnfAST) *ebnfNoneOrMore {
// 	return &ebnfNoneOrMore{child: child}
// }

// type ebnfOneOrMore struct {
// 	child ebnfAST
// }

// func (e *ebnfOneOrMore) match(s []rune) *token {
// 	lastMatchIdx := 0
// 	currIdx := 0

// 	tokens := make([]*token, 0)

// 	for {
// 		// go until token can't match or end of string
// 		for e.child.canMatch(s[lastMatchIdx:currIdx]) {
// 			currIdx++
// 			if currIdx > len(s) {
// 				break
// 			}
// 		}

// 		// back off one to possibly get valid token
// 		// Note it could be incomplete!
// 		currIdx--

// 		// create token
// 		token := e.child.match(s[lastMatchIdx:currIdx])
// 		if token == nil {
// 			return nil
// 		}
// 		tokens = append(tokens, token)

// 		if currIdx == len(s) {
// 			break
// 		}
// 		lastMatchIdx = currIdx
// 	}

// 	return &token{literal: s, tree: tokens}
// }

// func (e *ebnfOneOrMore) canMatch(s []rune) bool {
// 	lastMatchIdx := 0
// 	currIdx := 0

// 	for {
// 		// go until token can't match or end of string
// 		for e.child.canMatch(s[lastMatchIdx:currIdx]) {
// 			currIdx++
// 			if currIdx > len(s) {
// 				return true
// 			}
// 		}

// 		// back off one to possibly get valid token
// 		// Note it could be incomplete!
// 		currIdx--

// 		if e.child.match(s[lastMatchIdx:currIdx]) == nil {
// 			return false
// 		}

// 		if currIdx == len(s) {
// 			break
// 		}
// 		lastMatchIdx = currIdx
// 	}

// 	return true
// }

// func newOneOrMore(child ebnfAST) *ebnfOneOrMore {
// 	return &ebnfOneOrMore{child: child}
// }

// func main() {
// 	var tmp ebnfAST

// 	// [] - [&-&&&[&]]
// 	tmp = newException(newToRange(0, 1<<31-1), newCharRange("-&[]"))
// 	// rangeChar = ([] - [&-&&&[&]]) | "&&-" | "&&&&" | "&&[" | "&&]";
// 	rangeChar := newNonTerm("rangeChar", newAlt(tmp, newTerm("&-"), newTerm("&&"), newTerm("&["), newTerm("&]")))

// 	// identifierChar = [a-z] | [A-Z] | [0-9] | "_";
// 	identifierChar := newNonTerm("identifierChar", newAlt(newToRange('a', 'z'), newToRange('A', 'Z'), newTerm("_")))

// 	//[] - ["&&]
// 	tmp = newException(newToRange(0, 1<<31-1), newCharRange("\"&"))
// 	//stringChar = ([] - ["&&]) | "&&&"" | "&&&&";
// 	stringChar := newNonTerm("stringChar", newAlt(tmp, newTerm("&\""), newTerm("&&")))

// 	//spaceChar = [\t\n\v\f\rU+0085U+00A0];
// 	spaceChar := newCharRange("\t\n\v\f\rU+0085U+00A0 ")

// 	// terminal = "&"" stringChar* "&"";
// 	terminal := newNonTerm("terminal", newConcat(newTerm("\""), newNoneOrMore(stringChar), newTerm("\"")))

// 	//identifier = identifierChar+;
// 	identifier := newNonTerm("identifier", newOneOrMore(identifierChar))

// 	// toRange = (rangeChar "-" rangeChar);
// 	toRange := newNonTerm("toRange", newConcat(rangeChar, newTerm("-"), rangeChar))
// 	// charRange = rangeChar*;
// 	charRange := newNonTerm("charRange", newNoneOrMore(rangeChar))
// 	// range = "[" toRange | charRange "]";
// 	r := newNonTerm("range", newConcat(newTerm("["), newAlt(toRange, charRange), newTerm("]")))

// 	// space = spaceChar*;
// 	space := newNonTerm("space", newNoneOrMore(spaceChar))

// 	// we will come back to this!
// 	rhs := newNonTerm("RHS", nil)

// 	//term = "(" RHS ")"
// 	// | terminal
// 	// | identifier
// 	// | range;
// 	term := newNonTerm("term", newAlt(newConcat(newTerm("("), rhs, newTerm(")")), terminal, identifier, r))

// 	// ([?*+] | ("-" space term)
// 	tmp = newAlt(newCharRange("?*+"), newConcat(newTerm("-"), space, term))
// 	// factor = term (space ([?*+] | ("-" space term)))?;
// 	factor := newNonTerm("factor", newConcat(term, newNoneOrOne(newConcat(space, tmp))))

// 	//concatenation = factor (space factor)*;
// 	concatenation := newNonTerm("concatenation", newConcat(factor, newNoneOrMore(newConcat(space, factor))))

// 	//alternation = concatenation (space "|" space concatenation)*;
// 	alternation := newNonTerm("alternation", newConcat(concatenation, newNoneOrMore(newConcat(space, newTerm("|"), space, concatenation))))

// 	//RHS = space alternation space;
// 	rhs.rule = newConcat(space, alternation, space)

// 	//LHS = space identifier space;
// 	lhs := newNonTerm("lhs", newConcat(space, identifier, space))

// 	// rule = LHS "=" RHS ";";
// 	rule := newNonTerm("rule", newConcat(lhs, newTerm("="), rhs, newTerm(";")))

// 	//grammar = rule*;
// 	grammar := newNonTerm("grammar", newNoneOrMore(rule))



// 	data, err := os.ReadFile("test.txt")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	root := grammar.match([]rune(string(data)))
// 	root = findNonTerms(root)[0]
// 	print(root)
// } 

// func findNonTerms(root *token) []*token {
// 	newTree := make([]*token, 0)
// 	for _, tok := range root.tree {
// 		newTokens := findNonTerms(tok)
// 		newTree = append(newTree, newTokens...)
// 	}
// 	if root.nonTerm != "" {
// 		return []*token{{literal: root.literal, nonTerm: root.nonTerm, tree: newTree}}
// 	}
// 	return newTree
// }
