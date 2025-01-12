package grammar

import (
	"reflect"

	"asritha.dev/compiler/pkg/utils"
)

var (
	LongestRule = 0
)

type Rule struct {
	NonTerm        string
	SententialForm []*Symbol
}

func NewRule(nonTerm string, sententialForm ...*Symbol) *Rule {
	newRule := &Rule{
		NonTerm:        nonTerm,
		SententialForm: sententialForm,
	}

	len := len(newRule.String())
	if len > LongestRule {
		LongestRule = len
	}

	return newRule
}

func (r Rule) String() string {
	output := r.NonTerm + "="
	for _, s := range r.SententialForm {
		output += s.String()
	}
	return output
}

func (r Rule) Hash() int {
	return utils.HashStr(r.NonTerm) + utils.HashArr(r.SententialForm)
}

func (r Rule) Equal(other Rule) bool {
	return reflect.DeepEqual(r, other)
}
