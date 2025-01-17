package grammar

import (
	"reflect"

	"asritha.dev/compiler/pkg/utils"
)

type Rule struct {
	NonTerm        string
	SententialForm []Symbol
}

func NewRule(nonTerm string, sententialForm ...Symbol) Rule {
	return Rule{
		NonTerm:        nonTerm,
		SententialForm: sententialForm,
	}
}

/*
Returns the length of the sentential form not including epsilons
*/
func (r Rule) Len() int {
	len := 0
	for _, s := range r.SententialForm {
		if s != Epsilon {
			len++
		}
	}
	return len
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
