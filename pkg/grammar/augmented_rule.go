package grammar

import (
	"fmt"
	"asritha.dev/compiler/pkg/utils"
)

type simpleAugmentedRule struct {
	rule     *rule
	position int
}

func (ar simpleAugmentedRule) String() string {
	rule := ar.rule.nonTerm + "="
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

func (ar simpleAugmentedRule) StringWithLookahead(lookahead map[symbol]struct{}) string {
	rule := ar.rule.nonTerm + "="
	for i, s := range ar.rule.sententialForm {
		if ar.position == i {
			rule += "."
		}
		rule += s.String()
	}
	if ar.position == len(ar.rule.sententialForm) {
		rule += "."
	}

	return fmt.Sprintf("%-*s%v", longestRule + 4, rule, utils.MapToSetString(lookahead))
}

// type augmentedRule struct {
// 	simpleAugmentedRule
// 	lookahead map[symbol]struct{}
// }

// func NewAugmentedRule(r *rule, position int, lookahead map[symbol]struct{}) *augmentedRule {
// 	return &augmentedRule{
// 		simpleAugmentedRule: simpleAugmentedRule{r, position},
// 		lookahead:           lookahead,
// 	}
// }

// /*
// Returns a copy of the augment rule with the position shifted one to the right (next symbol)
// */
// func (ar augmentedRule) shiftedCopy() *augmentedRule {
// 	if ar.position == len(ar.rule.sententialForm) {
// 		return nil
// 	}

// 	return NewAugmentedRule(ar.rule, ar.position+1, ar.lookahead)
// }

// /*
// Returns the next symbol (symbol to right of position) in an augmented rule or nil if there is no next symbol
// */
// func (ar augmentedRule) getNextSymbol() *symbol {
// 	//Note: position should NOT be more than len(sententialForm)
// 	if ar.position == len(ar.rule.sententialForm) {
// 		return nil
// 	}
// 	return ar.rule.sententialForm[ar.position]
// }

// func (ar augmentedRule) String() string {
// 	rule := ar.rule.nonTerm + "="
// 	for i, s := range ar.rule.sententialForm {
// 		if ar.position == i {
// 			rule += "."
// 		}
// 		rule += s.String()
// 	}
// 	if ar.position == len(ar.rule.sententialForm) {
// 		rule += "."
// 	}

// 	return fmt.Sprintf("%-*s%v", longestRule + 4, rule, utils.MapToSetString(ar.lookahead))
// }

// func (ar augmentedRule) Hash() int {
// 	sum := 0
// 	for s := range ar.lookahead {
// 		sum += s.Hash()
// 	}
// 	return sum + ar.position + ar.rule.Hash()
// }

// func (ar augmentedRule) Equal(other augmentedRule) bool {
// 	return reflect.DeepEqual(ar, other)
// }
