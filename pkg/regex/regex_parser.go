package regex

import (
	"errors"
	"slices"
)

type regexParser struct {
	regex []rune
	curr  int // current index of unconsumed rune
}

// lookahead should just return the rune??? or is type used somewhere throughout
//func (rp *regexParser) lookAhead() rune {
//	return rp.regex[rp.curr]
//}

// is this called?
func (rp *regexParser) putBackToken() {
	rp.curr--
}

// TODO - replace this with above ?
// consumes and returns the current rune
//func (rp *regexParser) consume() rune {
//	if rp.curr > len(rp.regex) {
//		rp.curr++
//		return rp.regex[rp.curr-1]
//	}
//	return 0
//}

func (rp *regexParser) consumeIf(matching ...rune) bool {
	if rp.curr >= len(rp.regex) {
		return false // is this the best way to do ti?
	}
	if slices.Contains(matching, rp.regex[rp.curr]) {
		rp.curr++
		return true
	}
	return false
}

// Regex -> Alt
func (rp *regexParser) parse() (Node, error) {
	return rp.parseAlt()
}

// Alt -> Concat AltPrime
func (rp *regexParser) parseAlt() (Node, error) {
	node, err := rp.parseConcat()
	if err != nil {
		return nil, err
	}

	for {
		altPrime, err := rp.parseAltPrime()

		if err != nil {
			return nil, err
		}

		if altPrime != nil {
			node = NewAlternationNode(node, altPrime)
			continue
		}
		break
	}

	return node, nil
}

// AltPrime -> "|" Concat AltPrime | EPSILON
func (rp *regexParser) parseAltPrime() (Node, error) {
	if rp.consumeIf(firstSets["AltPrime"]...) {
		return rp.parseConcat()
	}

	return nil, nil
}

// Concat -> Repeat ConcatPrime
func (rp *regexParser) parseConcat() (Node, error) {
	node, err := rp.parseRepeat()
	if err != nil {
		return nil, err
	}

	for {
		concatPrime, err := rp.parseConcatPrime()

		if err != nil {
			return nil, err
		}

		if concatPrime != nil {
			node = NewConcatenationNode(concatPrime, node)
			continue
		}

		break
	}

	return node, nil
}

// ConcatPrime -> Repeat ConcatPrime | EPSILON
func (rp *regexParser) parseConcatPrime() (Node, error) {
	if rp.consumeIf(firstSets["ConcatPrime"]...) {
		return rp.parseRepeat()
	}
	return nil, nil
}

// Repeat -> Group Quantifier?
func (rp *regexParser) parseRepeat() (Node, error) {
	node, err := rp.parseGroup()
	if err != nil {
		return nil, err
	}

	if !rp.parsedQuantifier() { // there is no quantifier present
		return node, nil
	}

	quantifierToken := rp.regex[rp.curr-1]

	// some way to have the below in ParseQuantifier()?
	switch quantifierToken {
	case '*':
		// kleene star
		node = NewKleeneStarNode(node)
	case '+':
		// x+ can be rewritten as xx*
		node = NewConcatenationNode(node, NewKleeneStarNode(node))
	case '?':
		// x? can be rewritten as (s | EPSILON)
		node = NewAlternationNode(node, NewCharacterNode(0))
	}

	return node, nil
}

// Quantifier -> "*" | "+" | "?"
func (rp *regexParser) parsedQuantifier() bool {
	if rp.consumeIf(firstSets["Quantifier"]...) {
		return true
	}

	return false
}

// Group -> "(" Regex ")" | CharRange | Char
func (rp *regexParser) parseGroup() (Node, error) {
	var (
		node Node
		err  error
	)

	if rp.consumeIf('(') { // "(" Regex ")"
		node, err = rp.parse()
	} else if rp.consumeIf(firstSets["CharRange"]...) { // CharRange
		node, err = rp.parseCharRange()
	}

	if err != nil {
		return nil, err
	}

	return node, nil
}

// CharRange -> "[" CharRangeBody "]"
func (rp *regexParser) parseCharRange() (Node, error) {
	node, err := rp.parseCharRangeBody()
	if err != nil {
		return nil, err
	}

	rp.consumeIf(']')

	return node, nil
}

// CharRangeBody -> "^"? (CharRangeAtom)+
func (rp *regexParser) parseCharRangeBody() (Node, error) {
	isNot := false

	if rp.consumeIf('^') {
		isNot = true
	}

	node, err := rp.parseCharRangeAtom() // must consume CharRangeAtom at least once
	if err != nil {
		return nil, err
	}

	for slices.Contains(firstSets["CharRangeAtom"], rp.regex[rp.curr]) { // lookahead, better way of doing this?
		nextCharRangeAtom, err := rp.parseCharRangeAtom()
		if err != nil {
			return nil, err
		}
		node = NewAlternationNode(node, nextCharRangeAtom) // TODO: does this make sense?, should this be jsut concateanted of the CharacterClassNodes?
	}

	if !isNot { // if its true
		return node, nil
	}

	// invert

	isRunes := make([]rune, 0)
	isRunes = traverseCharRangeAtomOneOrMore(node, isRunes)

	notRunes := SubtractSlice2FromSlice1(anyChar, isRunes)

	return NewCharacterClassNode(notRunes), nil
}

// CharRangeAtom -> Char ("-" Char)?
func (rp *regexParser) parseCharRangeAtom() (Node, error) {
	startChar, err := rp.parseCharacter()
	if err != nil {
		return nil, err
	}

	if !rp.consumeIf('-') { // ("-" Char) is not present
		return startChar, nil
	}

	if !rp.consumeIf(firstSets["Char"]...) {
		return nil, errors.New("invalid character to the right of the -") // TODO: better error handling
	}

	indexOfStarChar := slices.Index(anyChar, rune(startChar))
	indexOfEndChar := slices.Index(anyChar, rp.regex[rp.curr-1]) // TODO: better way to do this?

	characterClass := anyChar[indexOfStarChar : indexOfEndChar+1]
	return NewCharacterClassNode(characterClass), nil
}

// Char -> ANY_VALID_CHAR
func (rp *regexParser) parseCharacter() (CharacterNode, error) {
	if rp.consumeIf(firstSets["Char"]...) {
		return NewCharacterNode(rp.regex[rp.curr-1]), nil // TODO: have a thing ofr rp.regex[rp.curr-1]??
	}
	return 0, errors.New("invalid char") // todo: fix error
}

// TODO: look at what interfaces the old regex parser implemented (String, ConvertToNFA)?
// TODO: is this best way to store ast?
// TODO: get rid of any nested ifs
// move to bottom, add comment
func traverseCharRangeAtomOneOrMore(node Node, runes []rune) []rune {
	if altNode, ok := node.(AlternationNode); ok {
		traverseCharRangeAtomOneOrMore(altNode.left, runes)
		traverseCharRangeAtomOneOrMore(altNode.right, runes)
	} else if characterNode, ok := node.(CharacterNode); ok {
		runes = append(runes, rune(characterNode))
	}
	return runes
}

// should be generic, move to utils, also does something like this alr exist?
// TODO: rename
func SubtractSlice2FromSlice1(slice1 []rune, slice2 []rune) []rune {
	newRunes := make([]rune, 0)
	//for _, i := range slice2 {
	//	if !slices.Contains(slice1, i) { // prolly way more efficient way to to do this
	//		newRunes = append(newRunes, i)
	//	}
	//}

	for _, i := range slice1 {
		if !slices.Contains(slice2, i) {
			newRunes = append(newRunes, i)
		}
	}
	return newRunes
}

// TODO: remove any ambiguity in consuming
