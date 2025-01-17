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

// type formatFunc func(*parseTreeNode) (*parseTreeNode, error)

// func GenerateRemoveFormatter(removeSet utils.Set[string]) formatFunc {
// 	return func(node *parseTreeNode) (*parseTreeNode, error) {
// 		if _, ok := removeSet[node.name]; ok {
// 			return nil, nil
// 		}
// 		return node, nil
// 	}
// }

// func GenerateRemoveEmptyFormatter() formatFunc {
// 	return func(node *parseTreeNode) (*parseTreeNode, error) {
// 		if node.literal == "" && len(node.children) == 0 {
// 			return nil, nil
// 		}
// 		return node, nil
// 	}
// }

// func GenerateCompressFormatter(compressSet utils.Set[string]) formatFunc {
// 	return func(node *parseTreeNode) (*parseTreeNode, error) {
// 		if _, ok := compressSet[node.name]; ok {
// 			node.literal = node.GetLiteral()
// 			node.children = node.children[:0]
// 		}
// 		return node, nil
// 	}
// }

/*
General purpose formatter. It starts by running the preChildrenFormatters using the output of the last formatter as the input to the next. It then recursively calls Format on all the children of the node then runs the postChildrenFormatters. If any formatter returns an error, it is brought up to the top level. If any formatter returns nil for a node, it is removed from its parent node.
*/
// func (node parseTreeNode) Format(preChildrenFormatters []formatFunc, postChildrenFormatters []formatFunc) (*parseTreeNode, error) {
// 	var err error
// 	for _, f := range preChildrenFormatters {
// 		node, err = f(node)
// 		if err != nil {
// 			return nil, err
// 		}
// 		if node == nil {
// 			return nil, nil
// 		}
// 	}

// 	for i := 0; i < len(node.children); i++ {
// 		node.children[i], err = node.children[i].Format(preChildrenFormatters, postChildrenFormatters)
// 		if err != nil {
// 			return nil, err
// 		}
// 		if node.children[i] == nil {
// 			utils.Remove(node.children, i)
// 			i--
// 		}
// 	}

// 	for _, f := range postChildrenFormatters {
// 		node, err = f(node)
// 		if err != nil {
// 			return nil, err
// 		}
// 		if node == nil {
// 			return nil, nil
// 		}
// 	}

// 	return node, nil
// }

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
