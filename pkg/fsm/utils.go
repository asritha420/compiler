package fsm

import (
	"fmt"
	"maps"
	"reflect"
	"slices"
)

/*
epsilonClosureAndTransitionsRecursion is the recursive method that will take an initial state and find the class containing all reachable states solely by epsilon. It will also get the ids and possible transitions of the class.
*/
func epsilonClosureAndTransitionsRecursion(initialState *NFAState, states []*NFAState, closed map[uint]struct{}, transitions map[rune]struct{}) ([]*NFAState) {
	if _, ok := closed[initialState.id]; ok {
		return states
	}
	closed[initialState.id] = struct{}{}
	states = append(states, initialState)
	for t := range maps.Keys(initialState.transitions) {
		transitions[t] = struct{}{}
	}
	for _, s := range initialState.transitions[Epsilon] {
		states = epsilonClosureAndTransitionsRecursion(s, states, closed, transitions)
	}
	return states
}

/*
epsilonClosureAndTransitions is a wrapper around epsilonClosureAndTransitionsRecursion which will get the class of states reachable solely by epsilon givin an initial class. It also returns the ids and possible transitions of the class.
*/
func epsilonClosureAndTransitions(states ...*NFAState) ([]*NFAState, map[uint]struct{}, map[rune]struct{}) {
	allStates := make([]*NFAState, 0)
	closed := make(map[uint]struct{})
	transitions := make(map[rune]struct{})
	for _, state := range states {
		allStates = epsilonClosureAndTransitionsRecursion(state, allStates, closed, transitions)
	}
	return allStates, closed, transitions
}

/*
isAccepting will check if any of the states in a class are accepting.
*/
func isAccepting(states ...*NFAState) bool {
	for _, state := range states {
		if state.isAccepting {
			return true
		}
	}
	return false
}

/*
idsToString will take a map of ids order them and convert them to a string. This ensures that 2 strings are identical iff the ids are identical.
*/
func idsToString[T any](ids map[uint]T) string {
	newIds := make([]uint, len(ids))
	for key := range ids {
		newIds = append(newIds, key)
	}
	slices.Sort(newIds)
	return fmt.Sprint(newIds)
}

type openListEntry struct {
	NFAstates   []*NFAState
	transitions map[rune]struct{}
	state       *NFAState
}

func ConvertNFAtoPseudoDFA(initialNFAState *NFAState) (*NFAState, map[string]*NFAState) {
	var id uint = 0
	DFAStates := make(map[string]*NFAState)
	openList := make([]openListEntry, 0)

	initialNFAClass, initialNFAClassIds, transitions := epsilonClosureAndTransitions(initialNFAState)
	
	initialDFAState := newPseudoDFAState(&id, isAccepting(initialNFAClass...))

	DFAStates[idsToString(initialNFAClassIds)] = initialDFAState
	openList = append(openList, openListEntry{
		NFAstates:   initialNFAClass,
		transitions: transitions,
		state:       initialDFAState,
	})

	for len(openList) > 0 {
		currentEntry := openList[0]
		openList = openList[1:]

		// loop through all possible transition (not including epsilon!)
		for transition := range currentEntry.transitions {
			if transition == Epsilon {
				continue
			}

			transitionNFAClass := make([]*NFAState, 0)
			// loop through all nodes in the current set and get all future nodes using the specific transition
			for _, currentNFAState := range currentEntry.NFAstates {
				transitionNFAClass = append(transitionNFAClass, currentNFAState.transitions[transition]...)
			}

			transitionNFAClass, transitionNFAIds, transitions := epsilonClosureAndTransitions(transitionNFAClass...)

			transitionNFAIdString := idsToString(transitionNFAIds)
			transitionDFAState, ok := DFAStates[transitionNFAIdString]
			// check if the transition leads to existing DFA
			if !ok {
				// if not make DFA state and add transition to open list
				transitionDFAState = newPseudoDFAState(&id, isAccepting(transitionNFAClass...))

				openList = append(openList, openListEntry{
					NFAstates:   transitionNFAClass,
					transitions: transitions,
					state:       transitionDFAState,
				})

				DFAStates[transitionNFAIdString] = transitionDFAState
			}

			// connect transition DFA to current DFA
			currentEntry.state.AddTransition(transition, transitionDFAState)
		}
	}

	return initialDFAState, DFAStates
}

type DFAClass struct {
	transitions map[rune]*DFAClass
	states      map[uint]*NFAState
	isAccepting bool
}

func findDFAClassFromState(state *NFAState, classes []*DFAClass) *DFAClass {
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
			transitions[key] = findDFAClassFromState(transition[0], classes)
		}
		if transitionClass := findDFAClassFromTransitions(transitions, new_classes); transitionClass != nil {
			transitionClass.states[state.id] = state
		} else {
			new_classes = append(new_classes, &DFAClass{
				transitions: transitions,
				states: map[uint]*NFAState{
					state.id: state,
				},
				isAccepting: class.isAccepting,
			})
		}
	}

	return new_classes
}

func MinimizeDFA(initialStateId uint, states []*NFAState) *NFAState {
	// Partition into accepting and non-accepting
	nonaccepting := &DFAClass{
		states:      make(map[uint]*NFAState),
		isAccepting: false,
	}
	accepting := &DFAClass{
		states:      make(map[uint]*NFAState),
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
			if len(seperatedClasses) != 1 {
				modified = true
			}
			newClasses = append(newClasses, seperatedClasses...)
		}
		classes = newClasses
	}

	// create all states
	classToState := make(map[string]*NFAState)
	var initialDFAState *NFAState
	for i, class := range classes {
		state := &NFAState{
			id:          uint(i),
			isAccepting: class.isAccepting,
			transitions: make(map[rune][]*NFAState),
			isPseudoDFA: true,
		}
		classToState[idsToString(class.states)] = state
		if _, ok := class.states[initialStateId]; ok {
			initialDFAState = state
		}
	}

	// link states
	for _, class := range classes {
		for transition, transitionClass := range class.transitions {
			classToState[idsToString(class.states)].AddTransition(transition, classToState[idsToString(transitionClass.states)])
		}
	}

	return initialDFAState
}
