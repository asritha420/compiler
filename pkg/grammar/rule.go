package grammar

import (
	"reflect"
	"asritha.dev/compiler/pkg/utils"
)

var (
	longestRule = 0
)

type rule struct {
	nonTerm        string
	sententialForm []*symbol
}

func NewRule(nonTerm string, sententialForm ...*symbol) *rule {
	newRule := &rule{
		nonTerm:        nonTerm,
		sententialForm: sententialForm,
	}
	
	len := len(newRule.String()) 
	if len > longestRule {
		longestRule = len 
	}

	return newRule
}

func (r rule) String() string {
	output := r.nonTerm + "="
	for _, s := range r.sententialForm {
		output += s.String()
	}
	return output
}

func (r rule) Hash() int {
	return utils.HashStr(r.nonTerm) + utils.HashArr(r.sententialForm)
}

func (r rule) Equal(other rule) bool {
	return reflect.DeepEqual(r, other)
}