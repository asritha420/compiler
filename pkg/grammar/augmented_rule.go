package grammar

import (
	"fmt"

	"asritha.dev/compiler/pkg/utils"
)

type augmentedRule struct {
	rule     *rule
	position int
}

func NewAugmentedRule(r *rule, position int) *augmentedRule {
	return &augmentedRule{
		rule:     r,
		position: position,
	}
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
	rule := ar.rule.NonTerm + "="
	for i, s := range ar.rule.sententialForm {
		if ar.position == i {
			rule += "."
		}
		rule += s.String()
	}
	if ar.position == len(ar.rule.sententialForm) {
		rule += "."
	}

	return rule
}

func (ar augmentedRule) StringWithLookahead(lookahead set[symbol]) string {
	rule := ar.rule.NonTerm + "="
	for i, s := range ar.rule.sententialForm {
		if ar.position == i {
			rule += "."
		}
		rule += s.String()
	}
	if ar.position == len(ar.rule.sententialForm) {
		rule += "."
	}

	return fmt.Sprintf("%-*s%v", longestRule+4, rule, utils.MapToSetString(lookahead))
}
