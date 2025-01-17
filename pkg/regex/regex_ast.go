package regex

type Node interface{}

type CharacterNode rune

func NewCharacterNode(char rune) CharacterNode {
	return CharacterNode(char)
}

type CharacterClassNode []rune

func NewCharacterClassNode(characterClass []rune) CharacterClassNode {
	return CharacterClassNode(characterClass)
}

type AlternationNode struct {
	left  Node
	right Node
}

func NewAlternationNode(left Node, right Node) AlternationNode {
	return AlternationNode{
		left:  left,
		right: right,
	}
}

type ConcatenationNode struct {
	left  Node
	right Node
}

func NewConcatenationNode(left Node, right Node) ConcatenationNode {
	return ConcatenationNode{
		left:  left,
		right: right,
	}
}

// Kleene start should not be a struct (it just has one field) TODO
type KleeneStarNode struct {
	left Node
}

func NewKleeneStarNode(left Node) KleeneStarNode {
	return KleeneStarNode{
		left: left,
	}
}
