package scanner

import (
	"asritha.dev/compiler/pkg/scanner/fsm"
	"fmt"
)

// TODO: in all of them there is no arrow pointing to the start state -> probably unnecessary tho
// separate out the types, printer, nfa converter into seperate files
const (
	epsilon rune = 0
)

// RExpr is implemented by all types automatically
type RExpr interface {
}

type ASTPrinter interface {
	PrintNode(indent string) string
}

// TODO: write algorithm name in comment here, comment explaining what it does?
type NFAConverter interface {
	convertToNFA(idCounter uint) (*scanner.NFAState, *scanner.NFAState) //start, end, create aliases?
}

type NFAPrinter interface {
	//printmermaidNFA()
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

// TODO: right now all of the IsAccepting is false, fix later
func (c *Const) convertToNFA(idCounter uint) (*scanner.NFAState, *scanner.NFAState) {
	endState := &scanner.NFAState{
		FAState: scanner.FAState{
			Id:          idCounter + 2,
			IsAccepting: false,
		},
		Transitions: make(map[rune][]*scanner.NFAState),
	}
	startState := &scanner.NFAState{
		FAState: scanner.FAState{
			Id:          idCounter + 1,
			IsAccepting: false,
		},
		Transitions: map[rune][]*scanner.NFAState{
			c.Value: []*scanner.NFAState{
				endState,
			},
		},
	}
	return startState, endState
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

func (a *Alternation) convertToNFA(idCounter uint) (*scanner.NFAState, *scanner.NFAState) {
	if left, ok := a.Left.(NFAConverter); ok {
		if right, ok := a.Right.(NFAConverter); ok {
			leftNFAStartState, leftNFAEndState := left.convertToNFA(idCounter)
			rightNFAStartState, rightNFAEndState := right.convertToNFA(idCounter)
			startState := &scanner.NFAState{
				FAState: scanner.FAState{
					Id:          idCounter + 1,
					IsAccepting: false,
				},
				Transitions: map[rune][]*scanner.NFAState{
					epsilon: []*scanner.NFAState{
						leftNFAStartState,
						rightNFAStartState,
					},
				},
			}
			endState := &scanner.NFAState{
				FAState: scanner.FAState{
					Id:          idCounter + 1,
					IsAccepting: true,
				},
				Transitions: make(map[rune][]*scanner.NFAState),
			}
			rightNFAEndState.IsAccepting = false
			leftNFAEndState.IsAccepting = false

			rightNFAEndState.Transitions[epsilon] = append(rightNFAEndState.Transitions[epsilon], endState)
			leftNFAEndState.Transitions[epsilon] = append(leftNFAEndState.Transitions[epsilon], endState)

			return startState, endState
		}
	}
	return nil, nil
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

func (c *Concatenation) convertToNFA(idCounter uint) (*scanner.NFAState, *scanner.NFAState) {
	if left, ok := c.Left.(NFAConverter); ok {
		if right, ok := c.Right.(NFAConverter); ok {
			//TODO: change the leftNFAEndState to not be final here, and in each corresponding function?
			leftNFAStartState, leftNFAEndState := left.convertToNFA(idCounter) //TODO: the idCounter is not being handled properly, return it everywhere and then increment by 1?
			rightNFAStartState, rightNFAEndState := right.convertToNFA(idCounter)
			leftNFAStartState.IsAccepting = false
			rightNFAEndState.IsAccepting = true
			leftNFAEndState.Transitions[epsilon] = append(leftNFAStartState.Transitions[epsilon], rightNFAStartState)

			return leftNFAStartState, rightNFAEndState
		} else {
			fmt.Println("Right part is invalid")
			//TODO: put in better error handling everywhere instead of just printing it
		}
	} else {
		fmt.Println("Left part is invalid")
	}
	return nil, nil
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

func (ks *KleeneStar) convertToNFA(idCounter uint) (*scanner.NFAState, *scanner.NFAState) {
	if left, ok := ks.Left.(NFAConverter); ok {
		leftNFAStartState, leftNFAEndState := left.convertToNFA(idCounter)

		leftNFAEndState.IsAccepting = false

		startState := &scanner.NFAState{
			FAState: scanner.FAState{
				Id:          idCounter + 1,
				IsAccepting: false,
			},
			Transitions: map[rune][]*scanner.NFAState{
				epsilon: []*scanner.NFAState{
					leftNFAStartState,
				},
			},
		}
		endState := &scanner.NFAState{
			FAState: scanner.FAState{
				Id:          idCounter + 1,
				IsAccepting: true,
			},
			Transitions: make(map[rune][]*scanner.NFAState),
		}

		//TODO: the below two prolly dont have to be separate, can blend together
		startState.Transitions[epsilon] = append(startState.Transitions[epsilon], endState)
		endState.Transitions[epsilon] = append(endState.Transitions[epsilon], startState)

		leftNFAEndState.Transitions[epsilon] = append(leftNFAEndState.Transitions[epsilon], endState)

		return startState, endState
	}
	return nil, nil
}

//TODO: make a mermaid js thing for this
