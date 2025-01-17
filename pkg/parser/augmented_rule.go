package parser

import (
	"fmt"

	"asritha.dev/compiler/pkg/grammar"
	"asritha.dev/compiler/pkg/utils"
)

type augmentedRule struct {
	*grammar.Rule
	position int
}

func NewAugmentedRule(r *grammar.Rule, position int) augmentedRule {
	//skip any epsilons
	for position < len(r.SententialForm) && r.SententialForm[position] == grammar.Epsilon {
		position++
	}
	
	return augmentedRule{
		Rule:     r,
		position: position,
	}
}

/*
Returns the next symbol (symbol to right of position) in an augmented rule or nil if there is no next symbol
*/
func (ar augmentedRule) getNextSymbol() grammar.Symbol {
	//Note: position should NOT be more than len(sententialForm)
	if ar.position == len(ar.SententialForm) {
		return grammar.Epsilon
	}
	return ar.SententialForm[ar.position]
}

func (ar augmentedRule) String() string {
	rule := ar.NonTerm + "="
	for i, s := range ar.SententialForm {
		if ar.position == i {
			rule += "."
		}
		rule += s.String()
	}
	if ar.position == len(ar.SententialForm) {
		rule += "."
	}

	return rule
}

func (ar augmentedRule) StringWithLookahead(lookahead utils.Set[grammar.Symbol], minSpacing int) string {
	rule := ar.NonTerm + " ="
	for i, s := range ar.SententialForm {
		if ar.position == i {
			rule += " ."
		}
		rule += " " + s.String()
	}
	if ar.position == len(ar.SententialForm) {
		rule += " ."
	}

	return fmt.Sprintf("%-*s%v", minSpacing, rule, utils.MapToSetString(lookahead))
}

// func (ar augmentedRule) Hash() int {
// 	return ar.Rule.Hash() + ar.position
// }

// func (ar augmentedRule) Equal(other augmentedRule) bool {
// 	return reflect.DeepEqual(ar, other)
// }