package grammar

import "reflect"

type augmentedRule struct {
	rule      *rule
	position  int
	lookahead map[symbol]struct{}
}

func NewAugmentedRule(r *rule, position int, lookahead map[symbol]struct{}) *augmentedRule {
	return &augmentedRule{
		rule:      r,
		position:  position,
		lookahead: lookahead,
	}
}

/*
Returns a copy of the augment rule with the position shifted one to the right (next symbol)
*/
func (ar augmentedRule) shiftedCopy(g *grammar) *augmentedRule {
	if ar.position == len(ar.rule.sententialForm) {
		return nil
	}
	return NewAugmentedRule(ar.rule, ar.position+1)
}

/*
Returns the next symbol (symbol to right of position) in an augmented rule or nil if there is no next symbol
*/
func (ar augmentedRule) getNextSymbol() *symbol {
	//Note: position should NOT be more than len(sententialForm)
	if ar.position == len(ar.rule.sententialForm) {
		return nil
	}
	return ar.rule.sententialForm[ar.position]
}

func (ar augmentedRule) String() string {
	// TODO add lookahead
	output := ar.rule.nonTerm + "="
	for i, s := range ar.rule.sententialForm {
		if ar.position == i {
			output += "."
		}
		output += s.String()
	}
	if ar.position == len(ar.rule.sententialForm) {
		output += "."
	}
	return output
}

func (ar augmentedRule) Hash() int {
	return ar.position + ar.rule.Hash()
}

func (ar augmentedRule) Equal(other augmentedRule) bool {
	return reflect.DeepEqual(ar, other)
}