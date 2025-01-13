package regex

func RegexParse() {

}

type RExpr interface{}

type Const struct {
	value rune
}

type Alternation struct {
	left  RExpr
	right RExpr
}

type Concatenation struct {
	left  RExpr
	right RExpr
}

type KleeneStar struct {
	left RExpr
}

// TODO: need this?
type RegexParser struct {
	tokens []*RegexToken
	curr   int // current index of unscanned token
}

func (rp *RegexParser) lookAhead() *RegexToken {
	return rp.tokens[rp.curr]
}

func (rp *RegexParser) putBackToken() {
	rp.curr--
}

func (rp *RegexParser) Parse() *RExpr {
	return rp.parseAlt()
}

func (rp *RegexParser) parseAlt() *RExpr {
	return rp.ParseConcat()
}

func (rp *RegexParser) ParseAltPrime() *RExpr {
	if rp.lookAhead().RegexTokenType == pipe {
		rp.curr++ // consume
		return &Alternation{
			left:
		}
	} else {

	}
}

func (rp *RegexParser) ParseConcat() *RExpr {
	return &Alternation{}
}

func (rp *RegexParser) ParseConcatPrime() *RExpr {}

//G// RExpr is implemented by all types automatically
//type RExpr interface {
//}
//
//type String interface {
//	String() string
//}
//
//// TODO add error handling?
//type ASTPrinter interface {
//	PrintNode(indent string) string
//}
//
//// Implements Thompson's Algorithm
//type NFAConverter interface {
//	convertToNFA(idCounter *uint) (*fsm.NFAState, *fsm.NFAState, error) //start, end, create aliases?
//}
//
//type Const struct {
//	value rune
//}
//
//func NewConst(value rune) *Const {
//	return &Const{value: value}
//}
//
//func (c *Const) String() string {
//	return fmt.Sprintf("%c", c.value)
//}
//
//func (c *Const) PrintNode(indent string) string {
//	return fmt.Sprintf("%sConst { %c }", indent, c.value)
//}
//
//func (c *Const) convertToNFA(idCounter *uint) (*fsm.NFAState, *fsm.NFAState, error) {
//	startState := fsm.NewNFAState(idCounter, false)
//	endState := fsm.NewNFAState(idCounter, true)
//	startState.AddTransition(c.value, endState)
//
//	return startState, endState, nil
//}
//
//type Alternation struct { // left | right
//	left  RExpr
//	right RExpr
//}
//
//func NewAlternation(left RExpr, right RExpr) *Alternation {
//	return &Alternation{left: left, right: right}
//}
//
//func (a *Alternation) String() string {
//	return fmt.Sprintf("%s|%s", a.left, a.right)
//}
//
//func (a *Alternation) PrintNode(indent string) string {
//	left, ok := a.left.(ASTPrinter)
//	if !ok {
//		return fmt.Sprintf("%sERROR PRINTING LEFT ALTERNATION", indent)
//	}
//
//	right, ok := a.right.(ASTPrinter)
//	if(!ok){
//		return fmt.Sprintf("%sERROR PRINTING RIGHT ALTERNATION", indent)
//	}
//
//	return fmt.Sprintf(
//		"%sAlternation {\n%v,\n%v\n%s}",
//		indent, left.PrintNode(indent+"  "),
//		right.PrintNode(indent+"  "),
//		indent,
//	)
//}
//
//// TODO add proper errors
//func (a *Alternation) convertToNFA(idCounter *uint) (*fsm.NFAState, *fsm.NFAState, error) {
//	left, ok := a.left.(NFAConverter)
//	if(!ok){
//		return nil, nil, fmt.Errorf("left fail")
//	}
//	right, ok := a.right.(NFAConverter)
//	if(!ok){
//		return nil, nil, fmt.Errorf("right fail")
//	}
//
//	leftNFAStartState, leftNFAEndState, err := left.convertToNFA(idCounter)
//	if err != nil {
//		return nil, nil, fmt.Errorf("left fail")
//	}
//	leftNFAEndState.IsAccepting = false
//	rightNFAStartState, rightNFAEndState, err := right.convertToNFA(idCounter)
//	if err != nil {
//		return nil, nil, fmt.Errorf("right fail")
//	}
//	rightNFAEndState.IsAccepting = false
//
//	startState := fsm.NewNFAState(idCounter, false)
//	endState := fsm.NewNFAState(idCounter, true)
//
//	startState.AddTransition(fsm.Epsilon, leftNFAStartState, rightNFAStartState)
//	leftNFAEndState.AddTransition(fsm.Epsilon, endState)
//	rightNFAEndState.AddTransition(fsm.Epsilon, endState)
//
//	return startState, endState, nil
//}
//
//type Concatenation struct { // left right
//	Left  RExpr
//	Right RExpr
//}
//
//func NewConcatenation(left RExpr, right RExpr) *Concatenation {
//	return &Concatenation{Left: left, Right: right}
//}
//
//func (c *Concatenation) String() string {
//	return fmt.Sprintf("%s%s", c.Left, c.Right)
//}
//
//func (c *Concatenation) PrintNode(indent string) string {
//	left, ok := c.Left.(ASTPrinter)
//	if !ok {
//		return fmt.Sprintf("%sERROR PRINTING LEFT CONCATENATION", indent)
//	}
//	right, ok := c.Right.(ASTPrinter)
//	if !ok {
//		return fmt.Sprintf("%sERROR PRINTING RIGHT CONCATENATION", indent)
//	}
//
//	return fmt.Sprintf(
//		"%sConcatenation {\n%v,\n%v\n%s}",
//		indent, left.PrintNode(indent+"  "),
//		right.PrintNode(indent+"  "),
//		indent,
//	)
//}
//
//func (c *Concatenation) convertToNFA(idCounter *uint) (*fsm.NFAState, *fsm.NFAState, error) {
//	left, ok := c.Left.(NFAConverter)
//	if(!ok){
//		return nil, nil, fmt.Errorf("left fail")
//	}
//	right, ok := c.Right.(NFAConverter)
//	if(!ok){
//		return nil, nil, fmt.Errorf("right fail")
//	}
//
//	leftNFAStartState, leftNFAEndState, err := left.convertToNFA(idCounter)
//	if err != nil {
//		return nil, nil, fmt.Errorf("left fail")
//	}
//	leftNFAEndState.IsAccepting = false
//	rightNFAStartState, rightNFAEndState, err := right.convertToNFA(idCounter)
//	if err != nil {
//		return nil, nil, fmt.Errorf("right fail")
//	}
//
//	leftNFAEndState.AddTransition(fsm.Epsilon, rightNFAStartState)
//
//	return leftNFAStartState, rightNFAEndState, nil
//}
//
//type KleeneStar struct { // left*
//	Left RExpr
//}
//
//func NewKleeneStar(left RExpr) *KleeneStar {
//	return &KleeneStar{Left: left}
//}
//
//func (ks *KleeneStar) String() string {
//	return fmt.Sprintf("(%s)*", ks.Left)
//}
//
//func (ks *KleeneStar) PrintNode(indent string) string {
//	if left, ok := ks.Left.(ASTPrinter); ok {
//		return fmt.Sprintf(
//			"%sKleeneStar {\n%v\n%s}",
//			indent, left.PrintNode(indent+"  "),
//			indent,
//		)
//	}
//	return fmt.Sprintf("%sERROR PRINTING KLEENE_STAR", indent)
//}
//
//func (ks *KleeneStar) convertToNFA(idCounter *uint) (*fsm.NFAState, *fsm.NFAState, error) {
//	left, ok := ks.Left.(NFAConverter)
//	if(!ok){
//		return nil, nil, fmt.Errorf("left fail")
//	}
//
//	childNFAStartState, childNFAEndState, err := left.convertToNFA(idCounter)
//	if err != nil {
//		return nil, nil, fmt.Errorf("left fail")
//	}
//	childNFAEndState.IsAccepting = false
//
//	startState := fsm.NewNFAState(idCounter, false)
//	endState := fsm.NewNFAState(idCounter, true)
//
//	startState.AddTransition(fsm.Epsilon, childNFAStartState, endState)
//	childNFAEndState.AddTransition(fsm.Epsilon, endState)
//	endState.AddTransition(fsm.Epsilon, startState)
//
//	return startState, endState, nil
//}import (
//"asritha.dev/compiler/pkg/scanner/internal/fsm"
//"fmt"
//)
//
//// RExpr is implemented by all types automatically
//type RExpr interface {
//}
//
//type String interface {
//	String() string
//}
//
//// TODO add error handling?
//type ASTPrinter interface {
//	PrintNode(indent string) string
//}
//
//// Implements Thompson's Algorithm
//type NFAConverter interface {
//	convertToNFA(idCounter *uint) (*fsm.NFAState, *fsm.NFAState, error) //start, end, create aliases?
//}
//
//type Const struct {
//	value rune
//}
//
//func NewConst(value rune) *Const {
//	return &Const{value: value}
//}
//
//func (c *Const) String() string {
//	return fmt.Sprintf("%c", c.value)
//}
//
//func (c *Const) PrintNode(indent string) string {
//	return fmt.Sprintf("%sConst { %c }", indent, c.value)
//}
//
//func (c *Const) convertToNFA(idCounter *uint) (*fsm.NFAState, *fsm.NFAState, error) {
//	startState := fsm.NewNFAState(idCounter, false)
//	endState := fsm.NewNFAState(idCounter, true)
//	startState.AddTransition(c.value, endState)
//
//	return startState, endState, nil
//}
//
//type Alternation struct { // left | right
//	left  RExpr
//	right RExpr
//}
//
//func NewAlternation(left RExpr, right RExpr) *Alternation {
//	return &Alternation{left: left, right: right}
//}
//
//func (a *Alternation) String() string {
//	return fmt.Sprintf("%s|%s", a.left, a.right)
//}
//
//func (a *Alternation) PrintNode(indent string) string {
//	left, ok := a.left.(ASTPrinter)
//	if !ok {
//		return fmt.Sprintf("%sERROR PRINTING LEFT ALTERNATION", indent)
//	}
//
//	right, ok := a.right.(ASTPrinter)
//	if(!ok){
//		return fmt.Sprintf("%sERROR PRINTING RIGHT ALTERNATION", indent)
//	}
//
//	return fmt.Sprintf(
//		"%sAlternation {\n%v,\n%v\n%s}",
//		indent, left.PrintNode(indent+"  "),
//		right.PrintNode(indent+"  "),
//		indent,
//	)
//}
//
//// TODO add proper errors
//func (a *Alternation) convertToNFA(idCounter *uint) (*fsm.NFAState, *fsm.NFAState, error) {
//	left, ok := a.left.(NFAConverter)
//	if(!ok){
//		return nil, nil, fmt.Errorf("left fail")
//	}
//	right, ok := a.right.(NFAConverter)
//	if(!ok){
//		return nil, nil, fmt.Errorf("right fail")
//	}
//
//	leftNFAStartState, leftNFAEndState, err := left.convertToNFA(idCounter)
//	if err != nil {
//		return nil, nil, fmt.Errorf("left fail")
//	}
//	leftNFAEndState.IsAccepting = false
//	rightNFAStartState, rightNFAEndState, err := right.convertToNFA(idCounter)
//	if err != nil {
//		return nil, nil, fmt.Errorf("right fail")
//	}
//	rightNFAEndState.IsAccepting = false
//
//	startState := fsm.NewNFAState(idCounter, false)
//	endState := fsm.NewNFAState(idCounter, true)
//
//	startState.AddTransition(fsm.Epsilon, leftNFAStartState, rightNFAStartState)
//	leftNFAEndState.AddTransition(fsm.Epsilon, endState)
//	rightNFAEndState.AddTransition(fsm.Epsilon, endState)
//
//	return startState, endState, nil
//}
//
//type Concatenation struct { // left right
//	Left  RExpr
//	Right RExpr
//}
//
//func NewConcatenation(left RExpr, right RExpr) *Concatenation {
//	return &Concatenation{Left: left, Right: right}
//}
//
//func (c *Concatenation) String() string {
//	return fmt.Sprintf("%s%s", c.Left, c.Right)
//}
//
//func (c *Concatenation) PrintNode(indent string) string {
//	left, ok := c.Left.(ASTPrinter)
//	if !ok {
//		return fmt.Sprintf("%sERROR PRINTING LEFT CONCATENATION", indent)
//	}
//	right, ok := c.Right.(ASTPrinter)
//	if !ok {
//		return fmt.Sprintf("%sERROR PRINTING RIGHT CONCATENATION", indent)
//	}
//
//	return fmt.Sprintf(
//		"%sConcatenation {\n%v,\n%v\n%s}",
//		indent, left.PrintNode(indent+"  "),
//		right.PrintNode(indent+"  "),
//		indent,
//	)
//}
//
//func (c *Concatenation) convertToNFA(idCounter *uint) (*fsm.NFAState, *fsm.NFAState, error) {
//	left, ok := c.Left.(NFAConverter)
//	if(!ok){
//		return nil, nil, fmt.Errorf("left fail")
//	}
//	right, ok := c.Right.(NFAConverter)
//	if(!ok){
//		return nil, nil, fmt.Errorf("right fail")
//	}
//
//	leftNFAStartState, leftNFAEndState, err := left.convertToNFA(idCounter)
//	if err != nil {
//		return nil, nil, fmt.Errorf("left fail")
//	}
//	leftNFAEndState.IsAccepting = false
//	rightNFAStartState, rightNFAEndState, err := right.convertToNFA(idCounter)
//	if err != nil {
//		return nil, nil, fmt.Errorf("right fail")
//	}
//
//	leftNFAEndState.AddTransition(fsm.Epsilon, rightNFAStartState)
//
//	return leftNFAStartState, rightNFAEndState, nil
//}
//
//type KleeneStar struct { // left*
//	Left RExpr
//}
//
//func NewKleeneStar(left RExpr) *KleeneStar {
//	return &KleeneStar{Left: left}
//}
//
//func (ks *KleeneStar) String() string {
//	return fmt.Sprintf("(%s)*", ks.Left)
//}
//
//func (ks *KleeneStar) PrintNode(indent string) string {
//	if left, ok := ks.Left.(ASTPrinter); ok {
//		return fmt.Sprintf(
//			"%sKleeneStar {\n%v\n%s}",
//			indent, left.PrintNode(indent+"  "),
//			indent,
//		)
//	}
//	return fmt.Sprintf("%sERROR PRINTING KLEENE_STAR", indent)
//}
//
//func (ks *KleeneStar) convertToNFA(idCounter *uint) (*fsm.NFAState, *fsm.NFAState, error) {
//	left, ok := ks.Left.(NFAConverter)
//	if(!ok){
//		return nil, nil, fmt.Errorf("left fail")
//	}
//
//	childNFAStartState, childNFAEndState, err := left.convertToNFA(idCounter)
//	if err != nil {
//		return nil, nil, fmt.Errorf("left fail")
//	}
//	childNFAEndState.IsAccepting = false
//
//	startState := fsm.NewNFAState(idCounter, false)
//	endState := fsm.NewNFAState(idCounter, true)
//
//	startState.AddTransition(fsm.Epsilon, childNFAStartState, endState)
//	childNFAEndState.AddTransition(fsm.Epsilon, endState)
//	endState.AddTransition(fsm.Epsilon, startState)
//
//	return startState, endState, nil
//}
