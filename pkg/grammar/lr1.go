package grammar

import (
	"fmt"
	"reflect"
	"strings"

	"asritha.dev/compiler/pkg/utils"
)

type lr1AutomationState struct {
	id             uint
	augmentedRules map[simpleAugmentedRule]set[symbol] //augmented rule -> lookahead
	transitions    map[symbol]*lr1AutomationState
}

func newLR1AutomationState(id *uint) *lr1AutomationState {
	*id++
	return &lr1AutomationState{
		id:             *id - 1,
		augmentedRules: make(map[simpleAugmentedRule]set[symbol]),
		transitions:    make(map[symbol]*lr1AutomationState),
	}
}

func (s *lr1AutomationState) String() string {
	rules := make([]string, len(s.augmentedRules))
	i := 0
	for ar, lookahead := range s.augmentedRules {
		rules[i] = ar.StringWithLookahead(lookahead)
		i++
	}
	return fmt.Sprintf("State %d\n%s", s.id, strings.Join(rules, "\n"))
}

func (ar simpleAugmentedRule) getClosureRecursion(g *Grammar, closure map[simpleAugmentedRule]set[symbol]) {
	nextSymbol := ar.getNextSymbol()
	if nextSymbol == nil || nextSymbol.symbolType != nonTerm {
		return
	}

	sentential := ar.rule.sententialForm[ar.position+1:]
	newLookahead := g.generateFirstSet(sentential...)

	if _, ok := newLookahead[Epsilon]; ok {
		utils.AddToMap(closure[ar], newLookahead)
	}

	delete(newLookahead, Epsilon)

	for _, rule := range g.ruleNTMap[nextSymbol.name] {
		newAR := *NewSimpleAugmentedRule(rule, 0)
		if _, ARExists := closure[newAR]; ARExists {
			utils.AddToMap(newLookahead, closure[newAR])
		} else {
			closure[newAR] = newLookahead
			newAR.getClosureRecursion(g, closure)
		}
	}
}

func getClosure(g *Grammar, initialClosure map[simpleAugmentedRule]set[symbol]) {
	for ar := range initialClosure {
		ar.getClosureRecursion(g, initialClosure)
	}
}

func (g *Grammar) getTransitions(core map[simpleAugmentedRule]set[symbol]) map[symbol]map[simpleAugmentedRule]set[symbol] {
	transitions := make(map[symbol]map[simpleAugmentedRule]set[symbol])

	for ar, lookahead := range core {
		nextSymbol := ar.getNextSymbol()
		if nextSymbol == nil {
			continue
		}

		if _, ok := transitions[*nextSymbol]; !ok {
			transitions[*nextSymbol] = make(map[simpleAugmentedRule]set[symbol])
		}

		newAR := *NewSimpleAugmentedRule(ar.rule, ar.position+1)
		//TODO may be redundant??
		if _, ok := transitions[*nextSymbol][newAR]; !ok {
			transitions[*nextSymbol][newAR] = make(set[symbol])
		}

		utils.AddToMap(lookahead, transitions[*nextSymbol][newAR])
		getClosure(g, transitions[*nextSymbol])
	}
	return transitions
}

func findState(target map[simpleAugmentedRule]set[symbol], states []*lr1AutomationState) *lr1AutomationState {
	for _, state := range states {
		if reflect.DeepEqual(target, state.augmentedRules) {
			return state
		}
	}
	return nil
}

func (g *Grammar) generateLR1() (*lr1AutomationState, []*lr1AutomationState) {
	var id uint = 0

	kernel := newLR1AutomationState(&id)
	kernel.augmentedRules = map[simpleAugmentedRule]set[symbol]{
		*NewSimpleAugmentedRule(g.Rules[0], 0): {EndOfInput:struct{}{}},
	}
	getClosure(g, kernel.augmentedRules)

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

	return kernel, states
}

func makeMermaid(states []*lr1AutomationState) string {
	mermaid := ""
	for _, state := range states {
		mermaid += fmt.Sprintf("state%d[\"`\n%s\n`\"]\n", state.id, state.String())
		for key, val := range state.transitions {
			transition := key.String()
			if transition == "+" {
				transition = "\\+"
			}
			mermaid += fmt.Sprintf("state%d--%s-->state%d\n", state.id, transition, val.id)
		}
	}
	return mermaid
}

// func convertLR1ToLALR(kernel *lr1AutomationState, states []*lr1AutomationState) (*lr1AutomationState, []*lr1AutomationState) {
// 	idMap := make(map[uint]uint)
// 	mergedStates :=
// }
