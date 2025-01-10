package grammar

import "testing"

func Test_rule_Hash(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		nonTerm        string
		sententialForm []*symbol
		want           int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRule(tt.nonTerm, tt.sententialForm)
			got := r.Hash()
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("Hash() = %v, want %v", got, tt.want)
			}
		})
	}
}
