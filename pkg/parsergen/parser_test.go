package parsergen

import "testing"

func TestTableParserGen_constructParseTable(t *testing.T) {
	tests := []struct {
		name string // description of this test case
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			var tp TableParserGen
			tp.constructParseTable()
		})
	}
}
