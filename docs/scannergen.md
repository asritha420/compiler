#### Introduction

This Go module will take in a text file containing regex rules and return an instance of a `Scanner` class.

#### Parsing Regex

The goal is to convert a regular expression $E$ into a parse tree.

We will use the following definition of a regular expression:

> A **regular expression** $s$ is a string which denotes $L(s)$, where $L(s)$ is the set of strings drawn from alphabet $\Sigma$.
>	$L(s)$ is defined inductively with the following base cases:
>		1) if $a \in \Sigma$, then $a$ is a regular expression and $L(a) = \{  a \}$.
>		2) if $\epsilon$ is a regular expression, then $L(\epsilon)$ contains only the empty string
>
Then for any regular expressions $s$ and $t$:
1) $s|t$ is an RE such that $L(s|t) = L(s) \cup L(t)$
2) $st$ is an RE such that $L(st)$ contains all the strings formed by concatenating a string in $L(s)$ followed by a string in $L(t)$
3) $s*$ is a RE such that $L(s*) = L(s)$ concatenated zero or more times

Using the above definition, here are a few examples of valid regular expressions: $ab|c$, $a(b|a)*a$, $(aaa)*$, etc.

We can start with the following grammar $G$:
- $E \rightarrow E \ | \ E$
- $E \rightarrow EE$ (*Note: represented as $E + E$ in the parse tree)
- $E \rightarrow E *$
- $E \rightarrow (E)$
- $E \rightarrow v$, where $v$ is any byte


The set of terminals $T$ is in $G$ is $T = \{ |, *, (, ), v\}$.

Note that $G$ is currently ambiguous, meaning it could lead to multiple valid parse tree derivations from the same input, and hence unpredictable, inconsistent behavior from the parser.
##### Removing Ambiguity in $G$

**Modifying $G$ to prevent Associativity Property Violation**:

We know the following:
> Alternation is left associative (ex. $a | b | c$ should be interpreted as $(a|b)|c$)
> Concatenation is left associative (ex. $abc$ should be interpreted as $(ab)c$)
> Kleene Closure is unary

Currently, $G$ violates the associativity property.

Take the following example regex: $s=abc$. Two possible parse trees can be produced from $G$ as of now, leading to ambiguity:

```
Possible Parse Tree 1: +(+(a, b), c)

	E 
  E +  E 
E + E  v
v   v  c
a   b 
```

```
Possible Parse Tree 2: +(a, +(b, c))

	E
 E  +   E 
 v    E + E
 a	  v   v
      b   c
```

We know from above only *Possible Parse Tree 1* is valid.

To satisfy left associativity, we can use left recursive rules; hence the grammar will expand the non-terminal $E$ from the leftmost occurrence, creating a parse tree that grows leftward, and forcing for operations to be evaluated from left to right.

Let's rewrite $G$ as the following:
- $E \rightarrow E \ | \ v$
- $E \rightarrow Ev$
- $E \rightarrow E *$
- $E \rightarrow (E)$
- $E \rightarrow v$, where $v$ is any byte


**Modifying $G$ to Prevent Precedence Property Violation**:

Continuing the definition of a regular expression, we know:

> Kleene closure ($s*$) the highest precedence.
> Concatenation ($st$) has next highest precedence.
> Alternation ($s|t$) has the lowest precedence.

Currently, $G$ violates the precedence property as per regex precedence rules.

Take the following example regex: $s = ab|c$, where $L(s) = \{ab, c \}$. Two possible parse trees can be produced from $G$ as of now, leading to ambiguity:

```
Possible Parse Tree 1: +(a, |(b,c))
	  E 
	E + E
	v   E | E
	a	v   v
		b   c
```


```
Possible Parse Tree 2: +(+(a, b), c)

	    E 
	E   | E
  E + E	  v
  v + v   c
  a   b
```


We know from the definition of a regex above that only *Possible Parse Tree 2* is valid.

As the generated parse tree is traversed in a "top to bottom," "left to right" manner, operators of higher precedence must be in the lower levels of the parse tree. #TODO: **mention top down parsing (leftmost derivation)**.

Let's rewrite $G$ as the following:

- $E \rightarrow [E\ | \ T] \ | \ T$ (*Note: $[, ]$ are NOT non terminals. They are simply used for grouping.)
- $T \rightarrow T F \ | \ F$
- $F \rightarrow G* | \ G$
- $G \rightarrow (E) \ | \ v$, where $v$ is any byte

Now the parse tree for $s$ would be:

```
Unambiguous Valid Parse Tree: +(+(a, b), c)

				  E 
				  T
			  T   +  F
		    T + F    G
			F   G    v
			G   v    c
			v   b 
			a
```

#TODO: check ambiguity
##### Eliminating Left Recursion in $G$

We must convert all our left recursive production rules ($A \rightarrow A\alpha | \beta$) to right recursive rules ($A \rightarrow \alpha A | \beta$), to avoid the possibility of an infinite loop during top down evaluation. #TODO: explain why left recursion is bad better

To eliminate left recursion, we can convert a left recursive production rule in the format $A \rightarrow A \alpha \ | \ \beta$ to $A \rightarrow \beta A'$ and $A' \rightarrow \alpha A ' \ | \ \epsilon$.

Let's rewrite $G$ as the following:

- $E \rightarrow TE'$
- $E' \rightarrow \ [\ | TE' \ ] \ | \ \epsilon$ (*Note: $[, ]$ are NOT terminals. They are simply used for grouping.)
- $T \rightarrow F T'$
- $T' \rightarrow FT' \ | \ \epsilon$
- $F \rightarrow G* | \ G$
- $G \rightarrow (E) \ | \ v$, where $v$ is any byte

##### Creating a Recursive Descent Parser for Regex

The parser will generate a parse tree for a given regex, if the regex is generated from $G$. The only decision the parser will make is what production to use.

Each non-terminal in the grammar will have a different function. 


package main

import (
"encoding/json"
"fmt"
"log"
"strings"
)

type RegexType int

const (
epsilon               = 0x00
Alternation RegexType = iota
Concatenation
KleeneStar
Constant
)

func main() {
// regexRules := map[string]string{
// 	"test": "a(cow|cat)*",
// }

	root := RegexConcatenationNode{
		LChildNode: RegexConstantNode{Data: 'a'},
		RChildNode: RegexKleeneStarNode{
			ChildNode: RegexAlternationNode{
				LChildNode: RegexConcatenationNode{
                    LChildNode: RegexConstantNode{Data: 'c'},
                    RChildNode: RegexConcatenationNode{
                        LChildNode: RegexConstantNode{Data: 'o'},
                        RChildNode: RegexConstantNode{Data: 'w'},
                    },
                },
				RChildNode: RegexConcatenationNode{
                    LChildNode: RegexConstantNode{Data: 'c'},
                    RChildNode: RegexConcatenationNode{
                        LChildNode: RegexConstantNode{Data: 'a'},
                        RChildNode: RegexConstantNode{Data: 't'},
                    },
                },
			},
		},
	}
	graph, _, _, err := convertRegexToNfa(root)
	if err != nil {
		log.Fatal(err)
	}
	DFA := convertNFAtoDFA(graph)
	// print(DFA)
	// println(m)
	println(makeMermaid(DFA))
}

type RegexASTNode interface{}

type RegexAlternationNode struct {
LChildNode RegexASTNode
RChildNode RegexASTNode
}

type RegexConcatenationNode struct {
LChildNode RegexASTNode
RChildNode RegexASTNode
}

type RegexKleeneStarNode struct {
ChildNode RegexASTNode
}

type RegexConstantNode struct {
Data byte
}

type StateEdge struct {
Transition byte
Next       State
}

type State interface {
GetId() uint
GetEdges() []StateEdge
IsAccepting() bool
}

type NFAState struct {
id          uint
transitions map[byte][]*NFAState //input:nfaState
isAccepting bool
}

func (state *NFAState) GetId() uint {
return state.id
}

func (state *NFAState) GetEdges() []StateEdge {
out := make([]StateEdge, 0)
for transition, nextStates := range state.transitions {
for _, nextState := range nextStates {
out = append(out, StateEdge{
Transition: transition,
Next:       nextState,
})
}
}
return out
}

func (state *NFAState) IsAccepting() bool {
return state.isAccepting
}

type DFAState struct {
id          uint
transitions map[byte]*DFAState
isAccepting     bool
}

func (state *DFAState) GetId() uint {
return state.id
}

func (state *DFAState) GetEdges() []StateEdge {
out := make([]StateEdge, 0)
for transitions, nextState := range state.transitions {
out = append(out, StateEdge{
Transition: transitions,
Next:       nextState,
})

	}
	return out
}

func (state *DFAState) IsAccepting() bool {
return state.isAccepting
}

func makeMermaidRecursion(rootState State, edges []string, closed map[uint]struct{}) ([]string, map[uint]struct{}) {
id := rootState.GetId()
if _, ok := closed[id]; ok {
return edges, closed
}
closed[id] = struct{}{}
for _, edge := range rootState.GetEdges() {
edges, closed = makeMermaidRecursion(edge.Next, edges, closed)
edges = append(edges, fmt.Sprintf("%d -- %d --> %d", id, edge.Transition, edge.Next.GetId()))

	}
	return edges, closed
}

func makeMermaid(rootState State) string {
edges, _ := makeMermaidRecursion(rootState, make([]string, 0), make(map[uint]struct{}))
return strings.Join(edges, "\n")
}

func convertRegexToNfaRecursion(regexASTRootNode RegexASTNode, idToState map[uint]*NFAState, id uint) (*NFAState, *NFAState, map[uint]*NFAState, uint, error) {

	switch rootNode := regexASTRootNode.(type) {
	case RegexConcatenationNode:
		lNFAState, lNFALastState, idToState, id, err := convertRegexToNfaRecursion(rootNode.LChildNode, idToState, id)
		if err != nil {
			return nil, nil, nil, 0, fmt.Errorf("concatination left node: %w", err)
		}
		rNFAState, rNFALastState, idToState, id, err := convertRegexToNfaRecursion(rootNode.RChildNode, idToState, id)
		if err != nil {
			return nil, nil, nil, 0, fmt.Errorf("concatenation right node: %w", err)
		}

		lNFALastState.transitions[epsilon] = append(lNFALastState.transitions[epsilon], rNFAState)
		return lNFAState, rNFALastState, idToState, id, nil
	case RegexAlternationNode:
		lNFAState, lNFALastState, idToState, id, err := convertRegexToNfaRecursion(rootNode.LChildNode, idToState, id)
		if err != nil {
			return nil, nil, nil, 0, fmt.Errorf("alternation left node: %w", err)
		}
		rNFAState, rNFALastState, idToState, id, err := convertRegexToNfaRecursion(rootNode.RChildNode, idToState, id)
		if err != nil {
			return nil, nil, nil, 0, fmt.Errorf("alternation right node: %w", err)
		}

		start := &NFAState{
			transitions: map[byte][]*NFAState{
				epsilon: {lNFAState, rNFAState},
			},
			id: id,
		}
		idToState[id] = start
		id++

		end := &NFAState{
			id:          id,
			transitions: make(map[byte][]*NFAState),
		}
		idToState[id] = end
		id++

		lNFALastState.transitions[epsilon] = append(lNFALastState.transitions[epsilon], end)
		rNFALastState.transitions[epsilon] = append(rNFALastState.transitions[epsilon], end)

		return start, end, idToState, id, nil
	case RegexKleeneStarNode:
		NFAStartState, NFALastState, idToState, id, err := convertRegexToNfaRecursion(rootNode.ChildNode, idToState, id)
		if err != nil {
			return nil, nil, nil, 0, fmt.Errorf("kleene star child node: %w", err)
		}

		start := &NFAState{
			transitions: map[byte][]*NFAState{
				epsilon: {NFAStartState},
			},
			id: id,
		}
		idToState[id] = start
		id++

		end := &NFAState{
			transitions: map[byte][]*NFAState{
				epsilon: {start},
			},
			id: id,
		}
		idToState[id] = end
		id++

		NFALastState.transitions[epsilon] = append(NFALastState.transitions[epsilon], end)
		start.transitions[epsilon] = append(start.transitions[epsilon], end)

		return start, end, idToState, id, nil
	case RegexConstantNode:
		start := &NFAState{
			transitions: make(map[byte][]*NFAState),
			id:          id,
		}
		idToState[id] = start
		id++

		end := &NFAState{
			id:          id,
			transitions: make(map[byte][]*NFAState),
		}
		idToState[id] = end
		id++

		start.transitions[rootNode.Data] = []*NFAState{end}
		return start, end, idToState, id, nil
	default:
		return nil, nil, nil, 0, fmt.Errorf("invalid node: %v", rootNode)
	}
}

func convertRegexToNfa(regexASTRootNode RegexASTNode) (*NFAState, *NFAState, map[uint]*NFAState, error) {
start, end, idMap, _, err := convertRegexToNfaRecursion(regexASTRootNode, make(map[uint]*NFAState), 0)
if err != nil {
return nil, nil, nil, fmt.Errorf("convertRegexToNfa: Unable to convert regex AST to NFA. Node trace:\n\t%w", err)
}
end.isAccepting = true
return start, end, idMap, nil
}

func epsilonClosureRecursion(initialState *NFAState, states []*NFAState, closed map[uint]struct{}) ([]*NFAState, map[uint]struct{}) {
if _, ok := closed[initialState.id]; ok {
return states, closed
}
closed[initialState.id] = struct{}{}
states = append(states, initialState)
for _, s := range initialState.transitions[epsilon] {
states, closed = epsilonClosureRecursion(s, states, closed)
}
return states, closed
}

func epsilonClosure(states ...*NFAState) ([]*NFAState, map[uint]struct{}) {
allStates := make([]*NFAState, 0)
closed := make(map[uint]struct{})
for _, state := range states {
allStates, closed = epsilonClosureRecursion(state, allStates, closed)
}
return allStates, closed
}

func IsAccepting(states ...*NFAState) bool {
for _, state := range states {
if state.IsAccepting() {
return true
}
}
return false
}

// TODO not use strings as ids??
func idsToString(ids map[uint]struct{}) string {
b, _ := json.Marshal(ids)
return string(b)
}

type OpenListEntry struct {
NFAstates []*NFAState
state     *DFAState
}

func convertNFAtoDFA(initialNFAState *NFAState) (*DFAState, map[string]*DFAState) {
var id uint = 0
DFAStates := make(map[string]*DFAState)
openList := make([]OpenListEntry, 0)

	initialNFAClass, initialNFAClassIds := epsilonClosure(initialNFAState)
	initialDFAState := &DFAState{
		id:          id,
		transitions: make(map[byte]*DFAState),
        isAccepting: IsAccepting(initialNFAClass...),
	}
	id++

	DFAStates[idsToString(initialNFAClassIds)] = initialDFAState
	openList = append(openList, OpenListEntry{
		NFAstates: initialNFAClass,
		state:     initialDFAState,
	})

	for len(openList) > 0 {
		currentEntry := openList[0]
		openList = openList[1:]

		// loop through all possible transition (not including epsilon!)
		for i := 1; i < 256; i++ {
			transition := byte(i)
			transitionNFAClass := make([]*NFAState, 0)
			// loop through all nodes in the current set and get all future nodes using the specific transition
			for _, currentNFAState := range currentEntry.NFAstates {
				transitionNFAClass = append(transitionNFAClass, currentNFAState.transitions[transition]...)
			}
			// if empty just continue (nothing will change)
			if len(transitionNFAClass) == 0 {
				continue
			}
			transitionNFAClass, transitionNFAIds := epsilonClosure(transitionNFAClass...)

			transitionNFAIdString := idsToString(transitionNFAIds)
			transitionDFAState, ok := DFAStates[transitionNFAIdString]
			// check if the transition leads to existing DFA
			if !ok {
				// if not make DFA state and add transition to open list
				transitionDFAState = &DFAState{
					id:          id,
					transitions: make(map[byte]*DFAState),
                    isAccepting: IsAccepting(transitionNFAClass...),
				}
				id++

				openList = append(openList, OpenListEntry{
					NFAstates: transitionNFAClass,
					state:     transitionDFAState,
				})

				DFAStates[transitionNFAIdString] = transitionDFAState
			}

			// connect transition DFA to current DFA
			currentEntry.state.transitions[transition] = transitionDFAState
		}
	}

	return initialDFAState, DFAStates
}

func minimizeDFA(initialDFAState *DFAState) *DFAState {
// Partition into accepting and non-accpeting

}