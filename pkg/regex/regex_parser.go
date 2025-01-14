package regex

import (
	"fmt"
	"slices"
)

var (
	anyValidChar = []rune{
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm',
		'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M',
		'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	}
	// TODO: FIX THIS MESS
	commonFirstSet        []rune  = append(anyValidChar, '(', '[')
	commonFirstSetPointer *[]rune = &commonFirstSet
	firstSets                     = map[string][]rune{
		"Regex":         *commonFirstSetPointer,
		"Alt":           *commonFirstSetPointer,
		"AltPrime":      {'|'},
		"Concat":        *commonFirstSetPointer,
		"ConcatPrime":   *commonFirstSetPointer,
		"Repeat":        *commonFirstSetPointer,
		"Quantifier":    {'*', '+', '?'},
		"Group":         *commonFirstSetPointer,
		"CharRange":     {'['},
		"CharRangeBody": append(anyValidChar, '^'),
		"CharRangeAtom": anyValidChar,
		"Char":          anyValidChar,
	}
)

type RExpr interface{}

type Const struct {
	value rune
}

func NewConst(value rune) Const {
	return Const{value}
}

type Alternation struct {
	left  RExpr
	right RExpr
}

func NewAlternation(left RExpr, right RExpr) Alternation {
	return Alternation{
		left:  left,
		right: right,
	}
}

type Concatenation struct {
	left  RExpr
	right RExpr
}

func NewConcatenation(left RExpr, right RExpr) Concatenation {
	return Concatenation{
		left:  left,
		right: right,
	}
}

type KleeneStar struct {
	left RExpr
}

func NewKleeneStar(left RExpr) KleeneStar {
	return KleeneStar{
		left: left,
	}
}

// TODO: need this?
type RegexParser struct {
	tokens []*RegexToken
	curr   int // current index of unscanned token
}

// TODO: add comments

// lookahead should just return the rune??? or is type used somewhere throughout
func (rp *RegexParser) lookAhead() *RegexToken {
	return rp.tokens[rp.curr]
}

func (rp *RegexParser) putBackToken() {
	rp.curr--
}

// Regex -> Alt
func (rp *RegexParser) ParseRegex() (RExpr, error) {
	return rp.parseAlt()
}

// Alt -> Concat AltPrime
func (rp *RegexParser) parseAlt() (RExpr, error) {
	expr, err := rp.parseConcat()
	if err != nil {
		return nil, err
	}

	for {
		altPrime, err := rp.parseAltPrime()

		if err != nil {
			return nil, err
		}

		if altPrime != nil {
			expr = NewAlternation(expr, altPrime)
			continue
		}
		break
	}

	return expr, nil
}

// AltPrime -> "|" Concat AltPrime | EPSILON
func (rp *RegexParser) parseAltPrime() (RExpr, error) {
	if slices.Contains(firstSets["AltPrime"], rp.lookAhead().rune) {
		rp.curr++
		return rp.parseConcat()
	}
	return nil, nil
}

// Concat -> Repeat ConcatPrime
func (rp *RegexParser) parseConcat() (RExpr, error) {
	expr, err := rp.parseRepeat()
	if err != nil {
		for {
			concatPrime, err := rp.parseConcatPrime()

			if err != nil {
				return nil, err
			}

			if concatPrime != nil {
				expr = NewConcatenation(concatPrime, expr)
				continue
			}

			break
		}
	}
	return expr, nil
}

// ConcatPrime -> Repeat ConcatPrime | EPSILON
func (rp *RegexParser) parseConcatPrime() (RExpr, error) {
	if slices.Contains(firstSets["ConcatPrime"], rp.lookAhead().rune) {
		rp.curr++
		return rp.parseRepeat()
	}
	return nil, nil
}

// Repeat -> Group Quantifier?
func (rp *RegexParser) parseRepeat() (RExpr, error) {
	expr, err := rp.parseGroup()
	if err != nil {
		return nil, err
	}

	if err = rp.parseQuantifier(); err != nil {
		return nil, err
	}

	quantifierToken := rp.tokens[rp.curr]
	switch quantifierToken.RegexTokenType {
	case star:
		// kleene star
		expr = NewKleeneStar(expr)
	case plus:
		// x+ can be rewritten as xx*
		expr = NewConcatenation(expr, NewKleeneStar(expr))
	case question:
		// x? can be rewritten as (s | EPSILON)
		expr = NewAlternation(expr, NewConst(0))
	default:
		// TODO: can it even hit this edge case?
		return nil, fmt.Errorf("ERROR") //TODO: fix error
	}

	return expr, nil
}

// Quantifier -> "*" | "+" | "?"
func (rp *RegexParser) parseQuantifier() error {
	if slices.Contains(firstSets["Quantifier"], rp.lookAhead().rune) {
		rp.curr++
	}
	return fmt.Errorf("error") // TODO: better error handling
}

// Group -> "(" Regex ")" | CharRange | Char
func (rp *RegexParser) parseGroup() (RExpr, error) {
	// TODO: need to do this?
	var (
		expr RExpr
		err  error
	)

	if rp.lookAhead().rune == '(' { // "(" Regex ")"
		rp.curr++
		expr, err = rp.ParseRegex()
		if err != nil {
			return nil, err
		}
	} else if slices.Contains(firstSets["CharRange"], rp.lookAhead().rune) { // CharRange
		expr, err = rp.parseCharRange()
	} else { // Char
		expr, err = rp.parseChar()
	}

	return expr, nil
}

// CharRange -> "[" CharRangeBody "]"
func (rp *RegexParser) parseCharRange() (RExpr, error) {
	rp.curr++ // consume "["

	expr, err := rp.parseCharRangeBody()
	if err != nil {
		return nil, err
	}

	rp.curr++ // consume "]"
	return expr, nil
}

// CharRangeBody -> "^"? (CharRangeAtom)+
func (rp *RegexParser) parseCharRangeBody() (RExpr, error) {
	not := false

	if rp.lookAhead().rune == '^' {
		not = true
		rp.curr++
	}

	expr, err := rp.parseCharRangeAtom() // consume it atleast once
	if err != nil {
		return nil, err
	}

	for !slices.Contains(firstSets["CharRangeAtom"], rp.lookAhead().rune) { // keep consuming
		nextCharRangeAtom, err := rp.parseCharRange()
		if err != nil {
			return nil, err
		}
		expr = NewAlternation(expr, nextCharRangeAtom)
	}

	// invert
	if not {
		// modify the expr to be not instead (invert it)
		// put this all in helper function?

		// traverse the expr tree via inorder
		isRunes := make([]rune, 0) // do the set ? ? do i need a set?

		var traverse func(root RExpr)
		traverse = func(root RExpr) {
			if altNode, ok := root.(Alternation); ok {
				traverse(altNode.left)
				traverse(altNode.right)
			} else if constNode, ok := root.(Const); ok {
				isRunes = append(isRunes, constNode.value)
			}
		}

		traverse(expr)

		// find the difference slice function
		isRunes = FindDifferenceSlices(anyValidChar, isRunes)

		var newExpr RExpr
		for i, r := range isRunes {
			if i == 0 {
				newExpr = NewConst(r)
			} else {
				newExpr = NewAlternation(newExpr, NewConst(r))
			}
		}

		return newExpr, nil
	}

	return expr, nil
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
func (rp *RegexParser) parseCharRangeAtom() (RExpr, error) {
	startChar, err := rp.parseChar() // rename this expr, or add comment saying its start char
	if err != nil {
		return nil, err
	}

	if rp.lookAhead().rune == '-' {
		rp.curr++ // consume '-'

		endChar := rp.lookAhead().rune

		indexOfStarChar := slices.Index(anyValidChar, startChar.(Const).value)

		indexOfEndChar := slices.Index(anyValidChar, endChar)

		for i := indexOfStarChar; i <= indexOfEndChar; i++ {
			startChar = NewAlternation(startChar, NewConst(anyValidChar[i]))
		}
	}

	return startChar, nil
}

// Char -> ANY_VALID_CHAR
func (rp *RegexParser) parseChar() (RExpr, error) {
	if slices.Contains(firstSets["Char"], rp.lookAhead().rune) {
		rp.curr++
		return NewConst(rp.tokens[rp.curr].rune), nil // create helper function for rp.tokens[rp.curr]?
	}
	return nil, fmt.Errorf("not a valid error") // TODO: fix
}

// TODO: look at what interfaces the old regex parser implemented (String, ConvertToNFA)?
// TODO: is this best way to store ast?
