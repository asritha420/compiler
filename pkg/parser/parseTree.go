package parser

import (
	"fmt"

	"asritha.dev/compiler/pkg/grammar"
	"asritha.dev/compiler/pkg/scanner"
)

type ParseTreeNode interface {
	GetLiteral() string
}

type ParseTreeNonTerm struct {
	rule     *grammar.Rule
	children []ParseTreeNode
}

func newParseTreeNonTerm(rule *grammar.Rule, children []ParseTreeNode) *ParseTreeNonTerm {
	return &ParseTreeNonTerm{
		rule:     rule,
		children: children,
	}
}

func (node ParseTreeNonTerm) GetLiteral() string {
	out := ""
	for _, c := range node.children {
		out += c.GetLiteral()
	}
	return out
}

type NonTermConversionFunc[T any] func(*ParseTreeNonTerm, []*T) (*T, error)
type TokenConversionFunc[T any] func(scanner.Token) (*T, error)

func Convert[T any](node ParseTreeNode, nonTermConversion map[*grammar.Rule]NonTermConversionFunc[T], nonTermDefaultConversion NonTermConversionFunc[T], tokenConversion TokenConversionFunc[T]) (*T, error) {
	switch n := node.(type) {
	case *ParseTreeNonTerm:
		children := make([]*T, len(n.children))
		for i, child := range n.children {
			ret, err := Convert(child, nonTermConversion, nonTermDefaultConversion, tokenConversion)
			if err != nil {
				return nil, err
			}
			children[i] = ret
		}

		if conversionFunc, ok := nonTermConversion[n.rule]; ok {
			return conversionFunc(n, children)
		}

		return nonTermDefaultConversion(n, children)
	case scanner.Token:
		return tokenConversion(n)
	}

	return nil, fmt.Errorf("invalid node %v", node)
}
