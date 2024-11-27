package scannergen

import (
	"fmt"
)

// RExpr is implemented by all types automatically
type RExpr interface{
	
}

type String interface {
	String() string
}

type ASTPrinter interface {
	PrintNode(indent string) string
}

type Const struct {
	Value rune
}

func NewConst(value rune) *Const {
	return &Const{Value: value}
}

func (c *Const) String() string {
	return fmt.Sprintf("%c", c.Value)
}

func (c *Const) PrintNode(indent string) string {
	return fmt.Sprintf("%sConst { %c }", indent, c.Value)
}

type Alternation struct { // left | right
	Left  RExpr
	Right RExpr
}

func NewAlternation(left RExpr, right RExpr) *Alternation {
	return &Alternation{Left: left, Right: right}
}

func (a *Alternation) String() string {
	return fmt.Sprintf("%s|%s", a.Left, a.Right)
}

func (a *Alternation) PrintNode(indent string) string {
	if left, ok := a.Left.(ASTPrinter); ok {
		if right, ok := a.Right.(ASTPrinter); ok {
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
	Left  RExpr
	Right RExpr
}

func NewConcatenation(left RExpr, right RExpr) *Concatenation {
	return &Concatenation{Left: left, Right: right}
}

func (c *Concatenation) String() string {
	return fmt.Sprintf("%s%s", c.Left, c.Right)
}

func (c *Concatenation) PrintNode(indent string) string {
	if left, ok := c.Left.(ASTPrinter); ok {
		if right, ok := c.Right.(ASTPrinter); ok {
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
	Left RExpr
}

func NewKleeneStar(left RExpr) *KleeneStar {
	return &KleeneStar{Left: left}
}

func (ks *KleeneStar) String() string {
	return fmt.Sprintf("(%s)*", ks.Left)
}

func (ks *KleeneStar) PrintNode(indent string) string {
	if left, ok := ks.Left.(ASTPrinter); ok {
		return fmt.Sprintf(
			"%sKleeneStar {\n%v\n%s}",
			indent, left.PrintNode(indent+"  "),
			indent,
		)
	}
	return fmt.Sprintf("%sERROR PRINTING KLEENE_STAR", indent)
}
