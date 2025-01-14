package parser

import (
	"asritha.dev/compiler/pkg/scannergenerator"
	"asritha.dev/compiler/pkg/utils"
)

type parseTreeNode struct {
	name     string
	literal  string
	children []*parseTreeNode
}

func newParseTreeNonTerm(name string, children []*parseTreeNode) *parseTreeNode {
	return &parseTreeNode{
		name:     name,
		children: children,
	}
}

func newParseTreeToken(t scannergenerator.Token) *parseTreeNode {
	return &parseTreeNode{
		name:    t.Name,
		literal: t.Literal,
	}
}

/*
Recursively constructs the literal represented by a node
*/
func (node *parseTreeNode) GetLiteral() string {
	if len(node.children) == 0 {
		return node.literal
	}
	out := ""
	for _, n := range node.children {
		out += n.GetLiteral()
	}
	return out
}


type formatFunc func(*parseTreeNode) (*parseTreeNode, error)

func GenerateRemoveFormatter(removeSet utils.Set[string]) formatFunc {
	return func(node *parseTreeNode) (*parseTreeNode, error) {
		if _, ok := removeSet[node.name]; ok {
			return nil, nil
		}
		return node, nil
	}
}

func GenerateRemoveEmptyFormatter() formatFunc {
	return func(node *parseTreeNode) (*parseTreeNode, error) {
		if node.literal == "" && len(node.children) == 0 {
			return nil, nil
		}
		return node, nil
	}
}

func GenerateCompressFormatter(compressSet utils.Set[string]) formatFunc {
	return func(node *parseTreeNode) (*parseTreeNode, error) {
		if _, ok := compressSet[node.name]; ok {
			node.literal = node.GetLiteral()
			node.children = node.children[:0]	
		}
		return node, nil
	}
}

/*
General purpose formatter. It starts by running the preChildrenFormatters using the output of the last one as the input to the next. It then recursively calls Format on all the children of the node then runs the postChildrenFormatters. If any formatter returns an error, it is brought up to the top level. If any formatter returns nil for a node, it is removed from its parent node.
*/
func (node *parseTreeNode) Format(preChildrenFormatters []formatFunc, postChildrenFormatters []formatFunc) (*parseTreeNode, error) {
	var err error
	for _, f := range preChildrenFormatters {
		node, err = f(node)
		if err != nil {
			return nil, err
		}
		if node == nil {
			return nil, nil
		}
	}

	for i := 0; i < len(node.children); i++ {
		node.children[i], err = node.children[i].Format(preChildrenFormatters, postChildrenFormatters)
		if err != nil {
			return nil, err
		}
		if node.children[i] == nil {
			utils.Remove(node.children, i)
			i--
		}
	}

	for _, f := range postChildrenFormatters {
		node, err = f(node)
		if err != nil {
			return nil, err
		}
		if node == nil {
			return nil, nil
		}
	}

	return node, nil
}

type convertFunc[T any] func(*parseTreeNode) (*T, error)
type convertFuncWithChildren[T any] func(*parseTreeNode, []*T) (*T, error)

func Convert[T any](node *parseTreeNode, preChildrenConversions map[string]convertFunc[T], postChildrenConversions map[string]convertFuncWithChildren[T], defaultConversion convertFuncWithChildren[T]) (*T, error) {
	if conversionFunc, ok := preChildrenConversions[node.name]; ok {
		return conversionFunc(node)
	}

	children := make([]*T, len(node.children))
	for i, child := range node.children {
		ret, err := Convert(child, preChildrenConversions, postChildrenConversions, defaultConversion)
		if err != nil {
			return nil, err
		}
		children[i] = ret
	}

	if conversionFunc, ok := postChildrenConversions[node.name]; ok {
		return conversionFunc(node, children)
	}

	return defaultConversion(node, children)
}