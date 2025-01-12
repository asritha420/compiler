package grammar

import (
	"reflect"

	"asritha.dev/compiler/pkg/utils"
)

type Rule struct {
	NonTerm        string
	SententialForm []*Symbol
}

func NewRule(nonTerm string, sententialForm ...*Symbol) *Rule {
	return &Rule{
		NonTerm:        nonTerm,
		SententialForm: sententialForm,
	}
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
