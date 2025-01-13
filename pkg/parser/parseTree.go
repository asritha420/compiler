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
Formats the parse tree by using the provided sets to remove extra nodes.

Note: order of ops goes removeSet > format children > compressSet > shorten
*/
func (node *parseTreeNode) Format(removeSet utils.Set[string], compressSet utils.Set[string], removeEmpty bool, shorten bool) *parseTreeNode {
	if _, ok := removeSet[node.name]; ok {
		return nil
	}

	for i := 0; i < len(node.children); i++ {
		child := node.children[i].Format(removeSet, compressSet, removeEmpty, shorten)
		if child == nil {
			node.children = utils.Remove(node.children, i)
			i--
		} else {
			node.children[i] = child
		}
	}

	if node.literal == "" && len(node.children) == 0 {
		return nil
	}

	if _, ok := compressSet[node.name]; ok {
		node.literal = node.GetLiteral()
		node.children = node.children[:0]
		return node
	}

	if shorten && len(node.children) == 1 {
		return node.children[0]
	}

	return node
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

func Convert[T any](node *parseTreeNode, conversionFuncs map[string]func(node *parseTreeNode, children []T) T, defaultFunc func(node *parseTreeNode, children []T) T ) T {
	children := make([]T, len(node.children))
	for i, child := range node.children {
		children[i] = Convert(child, conversionFuncs, defaultFunc)
	}

	if conversionFunc, ok := conversionFuncs[node.name]; ok {
		return conversionFunc(node, children)
	}

	return defaultFunc(node, children)
}