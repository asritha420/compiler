package scannergen

import (
	"fmt"
)

// RExpr is implemented by all types automatically
type RExpr interface{}

type String interface {
	String() string
}

type ASTPrinter interface {
	PrintNode(indent string) string
}

type Const struct {
	Value byte
}

func NewConst(value byte) *Const {
	return &Const{Value: value}
}

func (c *Const) String() string {
	return fmt.Sprintf("%c", c.Value)
}

func (c *Const) PrintNode(indent string) string {
	return fmt.Sprintf("%sConst { %c }", indent, c.Value)
}

type Alternation struct { // left | right
	left  RExpr
	right RExpr
}

func NewAlternation(left RExpr, right RExpr) *Alternation {
	return &Alternation{left: left, right: right}
}

func (a *Alternation) String() string {
	return fmt.Sprintf("%s|%s", a.left, a.right)
}

func (a *Alternation) PrintNode(indent string) string {
	if left, ok := a.left.(ASTPrinter); ok {
		if right, ok := a.right.(ASTPrinter); ok {
			return fmt.Sprintf(
				"%sAlternation {\n%v,\n%v\n%s}",
				indent, left.PrintNode(indent+"  "),
				right.PrintNode(indent+"  "),
				indent,
			)
		}
	}
	return fmt.Sprintf("%sERROR PRINTING ALTERNATION", indent)
}

type Concatenation struct { // left right
	left  RExpr
	right RExpr
}

func NewConcatenation(left RExpr, right RExpr) *Concatenation {
	return &Concatenation{left: left, right: right}
}

func (c *Concatenation) String() string {
	return fmt.Sprintf("%s%s", c.left, c.right)
}

func (c *Concatenation) PrintNode(indent string) string {
	if left, ok := c.left.(ASTPrinter); ok {
		if right, ok := c.right.(ASTPrinter); ok {
			return fmt.Sprintf(
				"%sConcatenation {\n%v,\n%v\n%s}",
				indent, left.PrintNode(indent+"  "),
				right.PrintNode(indent+"  "),
				indent,
			)
		}
	}
	return fmt.Sprintf("%sERROR PRINTING CONCATENATION", indent)
}

type KleeneStar struct { // left*
	left RExpr
}

func NewKleeneStar(left RExpr) *KleeneStar {
	return &KleeneStar{left: left}
}

func (ks *KleeneStar) String() string {
	return fmt.Sprintf("(%s)*", ks.left)
}

func (ks *KleeneStar) PrintNode(indent string) string {
	if left, ok := ks.left.(ASTPrinter); ok {
		return fmt.Sprintf(
			"%sKleeneStar {\n%v\n%s}",
			indent, left.PrintNode(indent+"  "),
			indent,
		)
	}
	return fmt.Sprintf("%sERROR PRINTING KLEENE_STAR", indent)
}
