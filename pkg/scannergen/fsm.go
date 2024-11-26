package scannergen

import (
	"encoding/json"
	"fmt"
	"maps"
	"reflect"
	"strings"
)

const (
	epsilon = 0x00
)

type StateEdge struct {
	Transition rune
	Next       State
}

type State interface {
	GetId() uint
	GetEdges() []StateEdge
	IsAccepting() bool
}

type NFAState struct {
	id          uint
	transitions map[rune][]*NFAState //input:nfaState
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
	transitions map[rune]*DFAState
	isAccepting bool
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

func getMermaidNodeIdString(state State) string {
	id := fmt.Sprintf("id%d", state.GetId())
	if state.IsAccepting() {
		id += fmt.Sprintf("(((%d)))", state.GetId())
	} else {
		id += fmt.Sprintf("((%d))", state.GetId())
	}
	return id
}

func makeMermaidRecursion(rootState State, edges []string, closed map[uint]struct{}) ([]string, map[uint]struct{}) {
	id := rootState.GetId()
	if _, ok := closed[id]; ok {
		return edges, closed
	}
	closed[id] = struct{}{}
	for _, edge := range rootState.GetEdges() {
		transition := ""
		if edge.Transition == epsilon {
			transition = "ɛ"
		} else {
			transition = fmt.Sprintf("%c", edge.Transition)
		}
		edges = append(edges, fmt.Sprintf("%s -- %s --> %s", getMermaidNodeIdString(rootState), transition, getMermaidNodeIdString(edge.Next)))
		edges, closed = makeMermaidRecursion(edge.Next, edges, closed)
	}
	return edges, closed
}

func MakeMermaid(rootState State) string {
	edges, _ := makeMermaidRecursion(rootState, make([]string, 0), make(map[uint]struct{}))
	edges = append(edges, fmt.Sprintf("START:::hidden -- start --> %s", getMermaidNodeIdString(rootState)))
	return strings.Join(edges, "\n")
}

func ConvertRegexToNfaRecursion(regexASTRootNode RExpr, idToState map[uint]*NFAState, id uint) (*NFAState, *NFAState, map[uint]*NFAState, uint, error) {
	switch rootNode := regexASTRootNode.(type) {
	case *Concatenation:
		lNFAState, lNFALastState, idToState, id, err := ConvertRegexToNfaRecursion(rootNode.Left, idToState, id)
		if err != nil {
			return nil, nil, nil, 0, fmt.Errorf("concatination left node: %w", err)
		}
		rNFAState, rNFALastState, idToState, id, err := ConvertRegexToNfaRecursion(rootNode.Right, idToState, id)
		if err != nil {
			return nil, nil, nil, 0, fmt.Errorf("concatenation right node: %w", err)
		}

		lNFALastState.transitions[epsilon] = append(lNFALastState.transitions[epsilon], rNFAState)
		return lNFAState, rNFALastState, idToState, id, nil
	case *Alternation:
		lNFAState, lNFALastState, idToState, id, err := ConvertRegexToNfaRecursion(rootNode.Left, idToState, id)
		if err != nil {
			return nil, nil, nil, 0, fmt.Errorf("alternation left node: %w", err)
		}
		rNFAState, rNFALastState, idToState, id, err := ConvertRegexToNfaRecursion(rootNode.Right, idToState, id)
		if err != nil {
			return nil, nil, nil, 0, fmt.Errorf("alternation right node: %w", err)
		}

		start := &NFAState{
			transitions: map[rune][]*NFAState{
				epsilon: {lNFAState, rNFAState},
			},
			id: id,
		}
		idToState[id] = start
		id++

		end := &NFAState{
			id:          id,
			transitions: make(map[rune][]*NFAState),
		}
		idToState[id] = end
		id++

		lNFALastState.transitions[epsilon] = append(lNFALastState.transitions[epsilon], end)
		rNFALastState.transitions[epsilon] = append(rNFALastState.transitions[epsilon], end)

		return start, end, idToState, id, nil
	case *KleeneStar:
		NFAStartState, NFALastState, idToState, id, err := ConvertRegexToNfaRecursion(rootNode.Left, idToState, id)
		if err != nil {
			return nil, nil, nil, 0, fmt.Errorf("kleene star child node: %w", err)
		}

		start := &NFAState{
			transitions: map[rune][]*NFAState{
				epsilon: {NFAStartState},
			},
			id: id,
		}
		idToState[id] = start
		id++

		end := &NFAState{
			transitions: map[rune][]*NFAState{
				epsilon: {start},
			},
			id: id,
		}
		idToState[id] = end
		id++

		NFALastState.transitions[epsilon] = append(NFALastState.transitions[epsilon], end)
		start.transitions[epsilon] = append(start.transitions[epsilon], end)

		return start, end, idToState, id, nil
	case *Const:
		start := &NFAState{
			transitions: make(map[rune][]*NFAState),
			id:          id,
		}
		idToState[id] = start
		id++

		end := &NFAState{
			id:          id,
			transitions: make(map[rune][]*NFAState),
		}
		idToState[id] = end
		id++

		start.transitions[rootNode.Value] = []*NFAState{end}
		return start, end, idToState, id, nil
	default:
		return nil, nil, nil, 0, fmt.Errorf("invalid node: %v", rootNode)
	}
}

func ConvertRegexToNfa(regexASTRootNode RExpr) (*NFAState, *NFAState, map[uint]*NFAState, error) {
	start, end, idMap, _, err := ConvertRegexToNfaRecursion(regexASTRootNode, make(map[uint]*NFAState), 0)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("convertRegexToNfa: Unable to convert regex AST to NFA. Node trace:\n\t%w", err)
	}
	end.isAccepting = true
	return start, end, idMap, nil
}

func epsilonClosureRecursion(initialState *NFAState, states []*NFAState, closed map[uint]struct{}, transitions map[rune]struct{}) ([]*NFAState, map[uint]struct{}, map[rune]struct{}) {
	if _, ok := closed[initialState.id]; ok {
		return states, closed, transitions
	}
	closed[initialState.id] = struct{}{}
	states = append(states, initialState)
	// TODO make more efficient?
	for t := range maps.Keys(initialState.transitions) {
		transitions[t] = struct{}{}
	}
	for _, s := range initialState.transitions[epsilon] {
		states, closed, transitions = epsilonClosureRecursion(s, states, closed, transitions)
	}
	return states, closed, transitions
}

func epsilonClosure(states ...*NFAState) ([]*NFAState, map[uint]struct{}, map[rune]struct{}) {
	allStates := make([]*NFAState, 0)
	closed := make(map[uint]struct{})
	transitions := make(map[rune]struct{})
	for _, state := range states {
		allStates, closed, transitions = epsilonClosureRecursion(state, allStates, closed, transitions)
	}
	return allStates, closed, transitions
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

func idsToStringDFA(ids map[uint]*DFAState) string {
	newMap := make(map[uint]struct{})
	for key, _ := range ids {
		newMap[key] = struct{}{}
	}
	return idsToString(newMap)
}

type OpenListEntry struct {
	NFAstates   []*NFAState
	Transitions map[rune]struct{}
	state       *DFAState
}

func ConvertNFAtoDFA(initialNFAState *NFAState) (*DFAState, map[string]*DFAState) {
	var id uint = 0
	DFAStates := make(map[string]*DFAState)
	openList := make([]OpenListEntry, 0)

	initialNFAClass, initialNFAClassIds, transitions := epsilonClosure(initialNFAState)
	initialDFAState := &DFAState{
		id:          id,
		transitions: make(map[rune]*DFAState),
		isAccepting: IsAccepting(initialNFAClass...),
	}
	id++

	DFAStates[idsToString(initialNFAClassIds)] = initialDFAState
	openList = append(openList, OpenListEntry{
		NFAstates:   initialNFAClass,
		Transitions: transitions,
		state:       initialDFAState,
	})

	for len(openList) > 0 {
		currentEntry := openList[0]
		openList = openList[1:]

		// loop through all possible transition (not including epsilon!)
		for transition := range currentEntry.Transitions {
			if transition == epsilon {
				continue
			}

			transitionNFAClass := make([]*NFAState, 0)
			// loop through all nodes in the current set and get all future nodes using the specific transition
			for _, currentNFAState := range currentEntry.NFAstates {
				transitionNFAClass = append(transitionNFAClass, currentNFAState.transitions[transition]...)
			}

			transitionNFAClass, transitionNFAIds, transitions := epsilonClosure(transitionNFAClass...)

			transitionNFAIdString := idsToString(transitionNFAIds)
			transitionDFAState, ok := DFAStates[transitionNFAIdString]
			// check if the transition leads to existing DFA
			if !ok {
				// if not make DFA state and add transition to open list
				transitionDFAState = &DFAState{
					id:          id,
					transitions: make(map[rune]*DFAState),
					isAccepting: IsAccepting(transitionNFAClass...),
				}
				id++

				openList = append(openList, OpenListEntry{
					NFAstates:   transitionNFAClass,
					Transitions: transitions,
					state:       transitionDFAState,
				})

				DFAStates[transitionNFAIdString] = transitionDFAState
			}

			// connect transition DFA to current DFA
			currentEntry.state.transitions[transition] = transitionDFAState
		}
	}

	return initialDFAState, DFAStates
}

type DFAClass struct {
	transitions map[rune]*DFAClass
	states      map[uint]*DFAState
	isAccepting bool
}

func findDFAClassFromState(state *DFAState, classes []*DFAClass) *DFAClass {
	for _, class := range classes {
		if _, ok := class.states[state.id]; ok {
			return class
		}
	}

	return nil
}

func findDFAClassFromTransitions(transitions map[rune]*DFAClass, classes []*DFAClass) *DFAClass {
	for _, class := range classes {
		if reflect.DeepEqual(class.transitions, transitions) {
			return class
		}
	}

	return nil
}

func makeDFAClasses(class *DFAClass, classes []*DFAClass) []*DFAClass {
	new_classes := make([]*DFAClass, 0)
	for _, state := range class.states {
		transitions := make(map[rune]*DFAClass)
		for key, transition := range state.transitions {
			transitions[key] = findDFAClassFromState(transition, classes)
		}
		if transitionClass := findDFAClassFromTransitions(transitions, new_classes); transitionClass != nil {
			transitionClass.states[state.id] = state
		} else {
			new_classes = append(new_classes, &DFAClass{
				transitions: transitions,
				states: map[uint]*DFAState{
					state.id: state,
				},
				isAccepting: class.isAccepting,
			})
		}
	}

	return new_classes
}

func MinimizeDFA(initialStateId uint, states map[string]*DFAState) *DFAState {
	// Partition into accepting and non-accepting
	nonaccepting := &DFAClass{
		states: make(map[uint]*DFAState),
		isAccepting: false,
	}
	accepting := &DFAClass{
		states: make(map[uint]*DFAState),
		isAccepting: true,
	}

	for _, state := range states {
		if state.isAccepting {
			accepting.states[state.id] = state
		} else {
			nonaccepting.states[state.id] = state
		}
	}

	accepting.isAccepting = true

	classes := []*DFAClass{
		accepting,
		nonaccepting,
	}

	modified := true
	for modified {
		modified = false
		newClasses := make([]*DFAClass, 0)
		for _, class := range classes {
			seperatedClasses := makeDFAClasses(class, classes)
			if(len(seperatedClasses) != 1) {
				modified = true	
			}
			newClasses = append(newClasses, seperatedClasses...)
		}
		classes = newClasses
	}

	// create all states
	classToState := make(map[string]*DFAState)
	var initialDFAState *DFAState
	for i, class := range classes {
		state := &DFAState{
			id: uint(i),
			isAccepting: class.isAccepting,
			transitions: make(map[rune]*DFAState),
		}
		classToState[idsToStringDFA(class.states)] = state
		if _, ok := class.states[initialStateId]; ok {
			initialDFAState = state
		}
	}

	// link states
	for _, class := range classes {
		for transition, transitionClass := range class.transitions {
			classToState[idsToStringDFA(class.states)].transitions[transition] = classToState[idsToStringDFA(transitionClass.states)]
		}
	}

	return initialDFAState
}
