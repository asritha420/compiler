package main

import (
	"reflect"
)

type ebnfAST interface { 
	matches([]rune) bool
}

type ebnfTerm struct {
	term []rune
}

func (t *ebnfTerm) matches(s []rune) bool {
	return reflect.DeepEqual(t.term, s)
}

type ebnfToRange struct {
	low, high rune
}

func (r *ebnfToRange) matches(s []rune) bool {
	if len(s) != 1 {
		return false
	}

	return s[0] >= r.low && s[0] < r.high
}

type ebnfCharRange struct {
	validChars []rune
}

func (r *ebnfCharRange) matches(s []rune) bool {
	if len(s) != 1 {
		return false
	}

	for _, c := range r.validChars {
		if s[0] == c {
			return true
		}
	}

	return false
}

type ebnfConcat struct {
	left, right ebnfAST
}

func (c *ebnfConcat) matches(s []rune) bool {
	for i := 0; i <= len(s); i++ {
		if c.left.matches(s[:i]) && c.right.matches(s[i:]) {
			return true
		}
	}
	return false
}

type ebnfAlt struct {
	left, right ebnfAST
}

func (a *ebnfAlt) matches(s []rune) bool {
	return a.left.matches(s) || a.right.matches(s)
}

type ebnfException struct {
	left, right ebnfAST
}

func (e *ebnfException) matches(s []rune) bool {
	return e.left.matches(s) && !e.right.matches(s)
}

type ebnfNoneOrOne struct {
	child ebnfAST
}

func (e *ebnfNoneOrOne) matches(s []rune) bool {
	return len(s) == 0 || e.child.matches(s)
}

type ebnfNoneOrMore struct {
	child ebnfAST
}

func (e *ebnfNoneOrMore) matches(s []rune) bool {
	if len(s) == 0 {
		return true
	}

	lastMatch := 0

	for i := 0; i <= len(s); i++ {
		if e.child.matches(s[lastMatch:i]) {
			lastMatch = i
		}
	}

	return lastMatch == len(s) 
}

type ebnfOneOrMore struct {
	child ebnfAST
}

func (e *ebnfOneOrMore) matches(s []rune) bool {
	if len(s) == 0 {
		return false
	}

	lastMatch := 0

	for i := 0; i <= len(s); i++ {
		if e.child.matches(s[lastMatch:i]) {
			lastMatch = i
		}
	}

	return lastMatch == len(s) 
}

func main() {
	rangeChar := &ebnfAlt{
		left: &ebnfException{
			left: &ebnfToRange{low: 0, high: 1<<30},
			right: &ebnfCharRange{
				validChars: []rune("-&"),
			},
		},
		right: &ebnfAlt{
			left: &ebnfTerm{term: []rune("&-")},
			right: &ebnfTerm{term: []rune("&&")},
		},
	}

	r := &ebnfConcat{
		left: &ebnfTerm{term: []rune("[")},
		right: &ebnfConcat{
			left: &ebnfAlt{
				left: &ebnfConcat{
					left: rangeChar,
					right: &ebnfConcat{
						left: &ebnfTerm{term: []rune("-")},
						right: rangeChar,
					},
				},
				right: &ebnfNoneOrMore{
					child: rangeChar,
				},
			},
			right: &ebnfTerm{term: []rune("]")},
		},
	}

	println(r.matches([]rune("[&&-hell]")))
}