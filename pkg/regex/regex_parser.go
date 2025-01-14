package regex

import (
	"errors"
	"fmt"
	"slices"
)

type regexParser struct {
	regex []rune
	curr  int // current index of unconsumed rune
}

// lookahead should just return the rune??? or is type used somewhere throughout
func (rp *regexParser) lookAhead() rune {
	return rp.regex[rp.curr]
}

// is this called?
func (rp *regexParser) putBackToken() {
	rp.curr--
}

// consumes and returns the current rune
func (rp *regexParser) consume() rune {
	if rp.curr > len(rp.regex) {
		rp.curr++
		return rp.regex[rp.curr-1]
	}
	return 0
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
	if slices.Contains(firstSets["AltPrime"], rp.lookAhead()) {
		rp.curr++ // consume "|"
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
	if slices.Contains(firstSets["ConcatPrime"], rp.lookAhead()) {
		rp.curr++
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

	if err = rp.parseQuantifier(); err != nil {
		return nil, err
	}

	quantifierToken := rp.regex[rp.curr]
	switch quantifierToken {
	case '*':
		// kleene star
		node = NewKleeneStarNode(node)
	case '+':
		// x+ can be rewritten as xx*
		node = NewConcatenationNode(node, NewKleeneStarNode(node))
	case '?':
		// x? can be rewritten as (s | EPSILON)
		node = NewAlternationNode(node, NewLiteralNode(0))
	}

	return node, nil
}

// Quantifier -> "*" | "+" | "?"
func (rp *regexParser) parseQuantifier() error {
	if slices.Contains(firstSets["Quantifier"], rp.lookAhead()) {
		rp.curr++ // consume quantifier
	}
	return fmt.Errorf("error") // TODO: better error handling
}

// Group -> "(" Regex ")" | CharRange | Char
func (rp *regexParser) parseGroup() (Node, error) {
	var (
		node Node
		err  error
	)

	if rp.lookAhead() == '(' { // "(" Regex ")"
		rp.curr++
		node, err = rp.parse()
	} else if slices.Contains(firstSets["CharRange"], rp.lookAhead()) { // CharRange
		node, err = rp.parseCharRange()
	} else { // Char
		node, err = rp.parseLiteral()
	}

	if err != nil {
		return nil, err
	}

	return node, nil
}

// CharRange -> "[" CharRangeBody "]"
func (rp *regexParser) parseCharRange() (Node, error) {
	rp.curr++ // consume "["

	node, err := rp.parseCharRangeBody()
	if err != nil {
		return nil, err
	}

	rp.curr++ // consume "]"
	return node, nil
}

// CharRangeBody -> "^"? (CharRangeAtom)+
func (rp *regexParser) parseCharRangeBody() (Node, error) {
	isNot := false

	if rp.lookAhead() == '^' {
		isNot = true
		rp.curr++ // consume "^"
	}

	node, err := rp.parseCharRangeAtom() // must consume CharRangeAtom at least once
	if err != nil {
		return nil, err
	}

	for !slices.Contains(firstSets["CharRangeAtom"], rp.lookAhead()) { // keep consuming CharRangeAtom
		nextCharRangeAtom, err := rp.parseCharRangeAtom()
		if err != nil {
			return nil, err
		}
		node = NewAlternationNode(node, nextCharRangeAtom) // TODO: does this make sense?, should this be jsut concateanted of the CharacterClassNodes?
	}

	if !isNot {
		return node, nil
	}

	// invert

	isRunes := make([]rune, 0)
	isRunes = traverseCharRangeAtomOneOrMore(node, isRunes)

	notRunes := FindDifferenceSlices(anyValidChar, isRunes)

	return NewCharacterClassNode(notRunes), nil
}

// move to bottom, add comment
func traverseCharRangeAtomOneOrMore(node Node, runes []rune) []rune {
	if altNode, ok := node.(AlternationNode); ok {
		traverseCharRangeAtomOneOrMore(altNode.left, runes)
		traverseCharRangeAtomOneOrMore(altNode.right, runes)
	} else if literalNode, ok := node.(LiteralNode); ok {
		runes = append(runes, rune(literalNode))
	}
	return runes
}

// should be generic, move to utils, also does something like this alr exist?
func FindDifferenceSlices(slice1 []rune, slice2 []rune) []rune {
	newRunes := make([]rune, 0)
	for _, i := range slice2 {
		if !slices.Contains(slice1, i) { // prolly way more efficient way to to do this
			newRunes = append(newRunes, i)
		}
	}
	return newRunes
}

// CharRangeAtom -> Char ("-" Char)?
func (rp *regexParser) parseCharRangeAtom() (Node, error) {
	startLiteral, err := rp.parseLiteral()
	if err != nil {
		return nil, err
	}

	if rp.lookAhead() != '-' { // ("-" Char) is not present
		return startLiteral, nil
	}

	rp.curr++ // consume '-'

	endChar := rp.consume()

	if endChar == 0 {
		return nil, errors.New(`nothing after "-"`)
	} // TODO: handle error, is 0 the correct one?

	indexOfStarChar := slices.Index(anyValidChar, rune(startLiteral))
	indexOfEndChar := slices.Index(anyValidChar, endChar) // TODO: better way to do this?

	characterClass := anyValidChar[indexOfStarChar : indexOfEndChar+1]
	return NewCharacterClassNode(characterClass), nil
}

// Char -> ANY_VALID_CHAR
func (rp *regexParser) parseLiteral() (LiteralNode, error) {
	if slices.Contains(firstSets["Char"], rp.lookAhead()) {
		rp.curr++ // consume the char
		return NewLiteralNode(rp.regex[rp.curr]), nil
	}
	return 0, errors.New("invalid char") // todo: fix error
}

// TODO: look at what interfaces the old regex parser implemented (String, ConvertToNFA)?
// TODO: is this best way to store ast?
// TODO: get rid of any nested ifs
