package grammar

import (
	"fmt"
	"reflect"
	"testing"
)

var (
	rules = []string{
		`production = expression`,
		`expression = term expressionPrime`,
		`expressionPrime = "|" term expressionPrime | " "`, //TODO: write in the spec how EPSILON should be specified as " "
		`term = factor termPrime`,
		`termPrime = factor termPrime | " "`,
		`factor = group factorPrime`,
		`factorPrime = "*" factorPrime | " "`,
		`group = "(" expression ")" | [a-z] | [A-Z] | [0-9]`,
	}

	nonTerminals = []string{"production", "expression", "term", "expressionPrime", "term", "termPrime", "factor", "group", "factorPrime"}
	terminals    = []string{"|", "*", "(", ")"} // TODO: write in comments or spec how the ranges don't have to be specified in the terminals array -> bc its confusing for the user, honestly might make more sense to just calculate it ourselves
)

// tests with regex grammar
func TestNewGrammar(t *testing.T) {
	result, err := NewGrammar(rules, nonTerminals, terminals)

	if err != nil {
		t.Errorf("NewGrammar returned an unexpected error: %s", err)
	}

	//TODO: could change implementation to use something like this ?
	regexGrammarSymbols := make(map[string]*symbol)
	// non-terminals
	for _, nT := range nonTerminals {
		regexGrammarSymbols[nT] = newNonTerminalSymbol(nT)
	}
	// terminals
	for _, t := range terminals {
		regexGrammarSymbols[t] = newTerminalSymbol(t)
	}
	// terminal ranges
	for _, terminalRange := range []symbolType{terminalLowercaseRange, terminalUppercaseRange, terminalNumberRange} {
		regexGrammarSymbols[terminalRange.String()] = newTerminalRangeSymbol(terminalRange)
	}

	// epsilon
	// TODO: define a newEpsilonSymbol()?
	regexGrammarSymbols["epsilon"] = &symbol{
		symbolType:    epsilon,
		validLiterals: []string{" "},
	}
	regexGrammar := &Grammar{
		Rules: []*Rule{
			{
				nonTerminal: "production",
				productions: []production{
					{
						regexGrammarSymbols["expression"],
					},
				},
			},
			{
				nonTerminal: "expression",
				productions: []production{
					{
						regexGrammarSymbols["term"],
						regexGrammarSymbols["expressionPrime"],
					},
				},
			},
			{
				nonTerminal: "expressionPrime",
				productions: []production{
					{
						regexGrammarSymbols["|"],
						regexGrammarSymbols["term"],
						regexGrammarSymbols["expressionPrime"],
					},
					{
						regexGrammarSymbols["epsilon"],
					},
				},
			},
			{
				nonTerminal: "term",
				productions: []production{
					{
						regexGrammarSymbols["factor"],
						regexGrammarSymbols["termPrime"],
					},
				},
			},
			{
				nonTerminal: "termPrime",
				productions: []production{
					{
						regexGrammarSymbols["factor"],
						regexGrammarSymbols["termPrime"],
					},
					{
						regexGrammarSymbols["epsilon"],
					},
				},
			},
			{
				nonTerminal: "factor",
				productions: []production{
					{
						regexGrammarSymbols["group"],
						regexGrammarSymbols["factorPrime"],
					},
				},
			},
			{
				nonTerminal: "factorPrime",
				productions: []production{
					{
						regexGrammarSymbols["*"],
						regexGrammarSymbols["factorPrime"],
					},
					{
						regexGrammarSymbols["epsilon"],
					},
				},
			},
			{
				nonTerminal: "group",
				productions: []production{
					{
						regexGrammarSymbols["("],
						regexGrammarSymbols["expression"],
						regexGrammarSymbols[")"],
					},
					{
						regexGrammarSymbols["terminalLowercaseRange"],
					},
					{
						regexGrammarSymbols["terminalUppercaseRange"],
					},
					{
						regexGrammarSymbols["terminalNumberRange"],
					},
				},
			},
		},
	}

	if !reflect.DeepEqual(result, regexGrammar) {
		t.Errorf("Result was incorrect: Got %s, want %s", result, regexGrammar)
	}
}

// TODO: move this to grammar file
func (g *Grammar) String() string {
	var rules string
	for _, r := range g.Rules {
		var productions string
		for _, p := range r.productions {
			var symbols string //TODO: should print symbolType literal instead of just number
			for _, s := range p {
				symbols = symbols + fmt.Sprintf("\n \t %+v", *s)
			} //TODO: should group each productions group together instead of just a new line separation
			productions = productions + symbols + "\n"
		}
		rules = rules + fmt.Sprintf("Rule {\n nonTerminal: \"%s\" \n productions: %s}, \n", r.nonTerminal, productions)
	}
	return rules
}

func TestGenerateFirstSet(t *testing.T) {
	g, err := NewGrammar(rules, nonTerminals, terminals)
	//TODO: repeating from above
	if err != nil {
		t.Errorf("TODO")
	}

	_ = [][]rune{
		{'(', 'v'}, // production
		{'(', 'v'}, // expression
		{},         // expressionPrime
		{'(', 'v'}, // term
		{},         // termPrime
		{'(', 'v'}, // factor
		{},         // factorPrime
		{'(', 'v'}, // group
	}

	g.generateFirstSets()

	/*
		rules = []string{
			`production = expression`,
			`expression = term expressionPrime`,
			`expressionPrime = "|" term expressionPrime | " "`, //TODO: write in the spec how EPSILON should be specified as " "
			`term = factor termPrime`,
			`termPrime = factor termPrime | " "`,
			`factor = group factorPrime`,
			`factorPrime = "*" factorPrime | " "`,
			`group = "(" expression ")" | [a-z] | [A-Z] | [0-9]`,
		}

			- start at first rule
			- for each production within each rule
			- loop through symbols in first rule
				- if it hits a nonTerminal, call the recursive getFirstSetForNonTerminal
					- set the returned results = to the first set of the rule's NT,
					- if it is not the last symbol in the production:
						- keep looping through rest of symbols. if it's a nonTerminal, call recursively. keep looping until it finds one w/ a first set that does not contain epsilon.
							- then add the first set of that to the og nonTerminal first set
					- if all of them have epsilon then add epsilon to the og first set
	*/

}
func TestGenerateFollowSet(t *testing.T) {}
