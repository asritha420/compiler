package grammar

import (
	"reflect"

	"asritha.dev/compiler/pkg/utils"
)

var (
	longestRule = 0
)

type rule struct {
	NonTerm        string
	sententialForm []*symbol
}

func NewRule(nonTerm string, sententialForm ...*symbol) *rule {
	newRule := &rule{
		NonTerm:        nonTerm,
		sententialForm: sententialForm,
	}

	len := len(newRule.String())
	if len > longestRule {
		longestRule = len
	}

	return newRule
}

func (r rule) String() string {
	output := r.NonTerm + "="
	for _, s := range r.sententialForm {
		output += s.String()
	}
	return output
}

func (r rule) Hash() int {
	return utils.HashStr(r.NonTerm) + utils.HashArr(r.sententialForm)
}

func (r rule) Equal(other rule) bool {
	return reflect.DeepEqual(r, other)
}
