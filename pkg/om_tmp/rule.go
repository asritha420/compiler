package grammar

import (
	"reflect"

	"asritha.dev/compiler/pkg/utils"
)

type rule struct {
	nonTerm        string
	sententialForm []*symbol
}

func NewRule(nonTerm string, sententialForm ...*symbol) *rule {
	return &rule{
		nonTerm:        nonTerm,
		sententialForm: sententialForm,
	}
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