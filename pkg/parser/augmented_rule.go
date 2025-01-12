package parser

import (
	"fmt"

	"asritha.dev/compiler/pkg/utils"
	."asritha.dev/compiler/pkg/grammar"
)

type augmentedRule struct {
	rule     *Rule
	position int
}

func NewAugmentedRule(r *Rule, position int) *augmentedRule {
	return &augmentedRule{
		rule:     r,
		position: position,
	}
}

/*
Returns the next symbol (symbol to right of position) in an augmented rule or nil if there is no next symbol
*/
func (ar augmentedRule) getNextSymbol() *Symbol {
	//Note: position should NOT be more than len(sententialForm)
	if ar.position == len(ar.rule.SententialForm) {
		return nil
	}
	return ar.rule.SententialForm[ar.position]
}

func (ar augmentedRule) String() string {
	rule := ar.rule.NonTerm + "="
	for i, s := range ar.rule.SententialForm {
		if ar.position == i {
			rule += "."
		}
		rule += s.String()
	}
	if ar.position == len(ar.rule.SententialForm) {
		rule += "."
	}

	return rule
}

func (ar augmentedRule) StringWithLookahead(lookahead utils.Set[Symbol], minSpacing int) string {
	rule := ar.rule.NonTerm + " ="
	for i, s := range ar.rule.SententialForm {
		if ar.position == i {
			rule += " ."
		}
		rule += " " + s.String()
	}
	if ar.position == len(ar.rule.SententialForm) {
		rule += " ."
	}

	return fmt.Sprintf("%-*s%v", minSpacing, rule, utils.MapToSetString(lookahead))
}
