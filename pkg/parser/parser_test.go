package parser_test

import(
	"asritha.dev/compiler/pkg/parser"
	"testing"
)

func TestNewSLRParser(t *testing.T) {
	tests := []struct {
		name string // description of this test case
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser.NewSLRParser()
		})
	}
}
