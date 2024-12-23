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
func epsilonClosureAndTransitionsRecursion(state *NFAState, states []*NFAState, closed map[uint]struct{}, transitions map[rune]struct{}) ([]*NFAState) {
	// check if we already found closure of state and add to closed list
	if _, ok := closed[state.id]; ok {
		return states
	}
	closed[state.id] = struct{}{}

	// add current state and transitions to output
	states = append(states, state)
	for t := range maps.Keys(state.transitions) {
		transitions[t] = struct{}{}
	}
	
	// recurse into all states from epsilon transition
	for _, s := range state.transitions[Epsilon] {
		states = epsilonClosureAndTransitionsRecursion(s, states, closed, transitions)
	}

	return states
}

/*
epsilonClosureAndTransitions is a wrapper around epsilonClosureAndTransitionsRecursion which will get the class of states reachable solely by epsilon givin an initial class. It also returns the ids and possible transitions of the class.
*/
func epsilonClosureAndTransitions(initialStates ...*NFAState) ([]*NFAState, map[uint]struct{}, map[rune]struct{}) {
	allStates := make([]*NFAState, 0)
	closed := make(map[uint]struct{})
	transitions := make(map[rune]struct{})
	// for all input states run epsilon closure
	for _, state := range initialStates {
		allStates = epsilonClosureAndTransitionsRecursion(state, allStates, closed, transitions)
	}
	return allStates, closed, transitions
}

/*
isAccepting will check if any of the states in a class are accepting.
*/
func isAccepting(states ...*NFAState) bool {
	for _, state := range states {
		if state.accepting {
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
	nfaStates   []*NFAState
	transitions map[rune]struct{}
	state       *NFAState
}

func ConvertNFAtoPseudoDFA(initialState *NFAState) (*NFAState, map[string]*NFAState) {
	var id uint = 0
	pseudoDFAStates := make(map[string]*NFAState)
	openList := make([]openListEntry, 0)

	initialNFASet, initialNFASetIds, initialNFASetTransitions := epsilonClosureAndTransitions(initialState)
	
	initialPseudoDFAState := NewNFAState(&id, isAccepting(initialNFASet...))

	pseudoDFAStates[idsToString(initialNFASetIds)] = initialPseudoDFAState
	openList = append(openList, openListEntry{
		nfaStates:   initialNFASet,
		transitions: initialNFASetTransitions,
		state:       initialPseudoDFAState,
	})

	for len(openList) > 0 {
		currentEntry := openList[0]
		openList = openList[1:]

		// loop through all possible transition (not including epsilon!)
		for transition := range currentEntry.transitions {
			if transition == Epsilon {
				continue
			}

			transitionNFASet := make([]*NFAState, 0)
			// loop through all nodes in the current set and get all future nodes using the specific transition
			for _, currentNFAState := range currentEntry.nfaStates {
				transitionNFASet = append(transitionNFASet, currentNFAState.transitions[transition]...)
			}

			transitionNFASet, transitionNFASetIds, transitionNFASetTransitions := epsilonClosureAndTransitions(transitionNFASet...)

			transitionNFAIdString := idsToString(transitionNFASetIds)
			transitionPseudoDFAState, ok := pseudoDFAStates[transitionNFAIdString]
			// check if the transition leads to existing DFA
			if !ok {
				// if not make DFA state and add transition to open list
				transitionPseudoDFAState = NewNFAState(&id, isAccepting(transitionNFASet...))

				openList = append(openList, openListEntry{
					nfaStates:   transitionNFASet,
					transitions: transitionNFASetTransitions,
					state:       transitionPseudoDFAState,
				})

				pseudoDFAStates[transitionNFAIdString] = transitionPseudoDFAState
			}

			// connect transition DFA to current DFA
			currentEntry.state.AddTransition(transition, transitionPseudoDFAState)
		}
	}

	return initialPseudoDFAState, pseudoDFAStates
}

type pseudoDFAClass struct {
	transitions map[rune]*pseudoDFAClass
	states      map[uint]*NFAState
	isAccepting bool
}

func findDFAClassFromState(state *NFAState, classes []*pseudoDFAClass) *pseudoDFAClass {
	for _, class := range classes {
		if _, ok := class.states[state.id]; ok {
			return class
		}
	}

	return nil
}

func findDFAClassFromTransitions(transitions map[rune]*pseudoDFAClass, classes []*pseudoDFAClass) *pseudoDFAClass {
	for _, class := range classes {
		if reflect.DeepEqual(class.transitions, transitions) {
			return class
		}
	}

	return nil
}

func makeDFAClasses(class *pseudoDFAClass, classes []*pseudoDFAClass) []*pseudoDFAClass {
	new_classes := make([]*pseudoDFAClass, 0)
	for _, state := range class.states {
		transitions := make(map[rune]*pseudoDFAClass)
		for key, transition := range state.transitions {
			transitions[key] = findDFAClassFromState(transition[0], classes)
		}
		if transitionClass := findDFAClassFromTransitions(transitions, new_classes); transitionClass != nil {
			transitionClass.states[state.id] = state
		} else {
			new_classes = append(new_classes, &pseudoDFAClass{
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

func MinimizePseudoDFA(initialStateId uint, states map[string]*NFAState) (*NFAState, map[string]*NFAState) {
	// Partition into accepting and non-accepting
	nonaccepting := &pseudoDFAClass{
		states:      make(map[uint]*NFAState),
		isAccepting: false,
	}
	accepting := &pseudoDFAClass{
		states:      make(map[uint]*NFAState),
		isAccepting: true,
	}

	for _, state := range states {
		if state.accepting {
			accepting.states[state.id] = state
		} else {
			nonaccepting.states[state.id] = state
		}
	}

	accepting.isAccepting = true

	classes := []*pseudoDFAClass{
		accepting,
		nonaccepting,
	}

	modified := true
	for modified {
		modified = false
		newClasses := make([]*pseudoDFAClass, 0)
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
	var i uint = 0
	for _, class := range classes {
		state := NewNFAState(&i, class.isAccepting)
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

	return initialDFAState, classToState
}
