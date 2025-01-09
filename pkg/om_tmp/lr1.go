package grammar

import "fmt"

type lr1AutomationState struct {
	id             uint
	augmentedRules map[*augmentedRule]struct{}
	transitions    map[symbol]*lr1AutomationState
}

func newLR1AutomationState(id *uint) *lr1AutomationState {
	*id++
	return &lr1AutomationState{
		id:             *id - 1,
		augmentedRules: make(map[*augmentedRule]struct{}),
		transitions:    make(map[symbol]*lr1AutomationState),
	}
}

func (g *grammar) getClosureRecursion(ar *augmentedRule, closure map[*augmentedRule]struct{}, closed map[string]struct{}) {
	if _, ok := closed[ar.String()]; ok {
		return
	}

	closed[ar.String()] = struct{}{}
	closure[ar] = struct{}{}

	nextSymbol := ar.getNextSymbol()
	if nextSymbol == nil || nextSymbol.sType != NonTerm {
		return
	}

	// nextSymbol is a NT
	for _, r := range g.ruleNTMap[nextSymbol.data] {
		newAR := NewAugmentedRule(r, 0)
		g.getClosureRecursion(newAR, closure, closed)
	}
}

func (g *grammar) getClosure(ars ...*augmentedRule) map[*augmentedRule]struct{} {
	closure := make(map[*augmentedRule]struct{})
	closed := make(map[string]struct{})

	for _, ar := range ars {
		g.getClosureRecursion(ar, closure, closed)
	}

	return closure
}

func (g *grammar) getTransitions(ars map[*augmentedRule]struct{}) map[symbol]map[*augmentedRule]struct{} {
	transitions := make(map[symbol]map[*augmentedRule]struct{})
	closed := make(map[symbol]map[string]struct{})
	for ar := range ars {
		nextSymbol := ar.getNextSymbol()
		if nextSymbol == nil {
			continue
		}

		if _, ok := transitions[*nextSymbol]; !ok {
			transitions[*nextSymbol] = make(map[*augmentedRule]struct{})
			closed[*nextSymbol] = make(map[string]struct{})
		}
		g.getClosureRecursion(ar.shiftedCopy(), transitions[*nextSymbol], closed[*nextSymbol])
	}
	return transitions
}

func equal(m1, m2 map[*augmentedRule]struct{}) bool {
	if len(m1) != len(m2) {
		return false
	}

	strs := make(map[string]struct{})
	for ar := range m1 {
		strs[ar.String()] = struct{}{}
	}

	for ar := range m2 {
		if _, ok := strs[ar.String()]; !ok {
			return false
		}
	}

	return true
}

func findState(target map[*augmentedRule]struct{}, states []*lr1AutomationState) *lr1AutomationState {
	for _, state := range states {
		fmt.Printf("target: %v\n", target)
		fmt.Printf("state : %v\n", state.augmentedRules)
		if equal(target, state.augmentedRules) {
			return state
		}
	}
	return nil
}

func (g *grammar) generateLR1() *lr1AutomationState {
	var id uint = 0

	kernel := newLR1AutomationState(&id)
	startRule := NewAugmentedRule(g.rules[0], 0)
	kernel.augmentedRules = g.getClosure(startRule)

	states := []*lr1AutomationState{kernel}
	openList := []*lr1AutomationState{kernel}

	for len(openList) > 0 {
		state := openList[0]
		openList = openList[1:]

		for transition, rules := range g.getTransitions(state.augmentedRules) {
			if transitionState := findState(rules, states); transitionState != nil {
				state.transitions[transition] = transitionState
				continue
			}

			newState := newLR1AutomationState(&id)
			newState.augmentedRules = rules
			state.transitions[transition] = newState
			states = append(states, newState)
			openList = append(openList, newState)
		}
	}

	return kernel
}