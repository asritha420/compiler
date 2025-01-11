package grammar

import "testing"

func Test_getClosure(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		g       *Grammar
		initial map[simpleAugmentedRule]set[symbol]
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getClosure(tt.g, tt.initial)
		})
	}
}

func Test_simpleAugmentedRule_getClosureRecursion(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		r        *rule
		position int
		// Named input parameters for target function.
		g       *Grammar
		closure map[simpleAugmentedRule]set[symbol]
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ar := NewSimpleAugmentedRule(tt.r, tt.position)
			ar.getClosureRecursion(tt.g, tt.closure)
		})
	}
}
