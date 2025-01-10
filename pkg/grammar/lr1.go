package grammar

import (
	"fmt"
	"strings"

	"asritha.dev/compiler/pkg/utils"
)

type lr1AutomationState struct {
	id             uint
	augmentedRules *utils.Map[augmentedRule, struct{}]
	transitions    map[symbol]*lr1AutomationState
}

func newLR1AutomationState(id *uint) *lr1AutomationState {
	*id++
	return &lr1AutomationState{
		id:             *id - 1,
		augmentedRules: utils.NewMap[augmentedRule, struct{}](),
		transitions:    make(map[symbol]*lr1AutomationState),
	}
}

func (s *lr1AutomationState) String() string {
	rules := make([]string,s.augmentedRules.Len())
	for i, ar := range s.augmentedRules.GetAllKeys() {
		rules[i] = ar.String()
	}
	return fmt.Sprintf("State %d\n%s", s.id, strings.Join(rules, "\n"))
}

func (ar *augmentedRule) getClosureRecursion(g *grammar, closure *utils.Map[augmentedRule, struct{}]) {
	if _, ok := closure.Get(*ar); ok {
		return
	}

	closure.Put(*ar, struct{}{})

	nextSymbol := ar.getNextSymbol()
	if nextSymbol == nil || nextSymbol.symbolType != nonTerm {
		return
	}

	// nextSymbol is a NT
	newLookahead := make(map[symbol]struct{})
	sentential := ar.rule.sententialForm[ar.position + 1:]
	first := g.generateFirstSet(sentential...)
	utils.AddToMap(first, newLookahead)
	if _, containsEpsilon := first[Epsilon]; containsEpsilon {
		utils.AddToMap(ar.lookahead, newLookahead)
	}

	delete(newLookahead, Epsilon)

	
	for _, r := range g.ruleNTMap[nextSymbol.name] {
		newAR := NewAugmentedRule(r, 0, newLookahead)
		newAR.getClosureRecursion(g, closure)
	}
}

//TODO find a better way to do this?
func getClosure(g *grammar, ars ...*augmentedRule) *utils.Map[augmentedRule,struct{}] {
	closure := utils.NewMap[augmentedRule, struct{}]()

	for _, ar := range ars {
		ar.getClosureRecursion(g, closure)
	}

	// mergedClosure := utils.NewMap[augmentedRule, struct{}]()
	arLookaheadMap := make(map[simpleAugmentedRule]map[symbol]struct{})
	for _, ar := range closure.GetAllKeys() {
		if lookahead, ok := arLookaheadMap[ar.simpleAugmentedRule]; ok {
			closure.Remove(ar)
			utils.AddToMap(ar.lookahead, lookahead)
		} else {
			arLookaheadMap[ar.simpleAugmentedRule] = ar.lookahead
		}
	}

	return closure
}

func (g *grammar) getTransitions(ars *utils.Map[augmentedRule, struct{}]) map[symbol]*utils.Map[augmentedRule, struct{}] {
	transitions := make(map[symbol]*utils.Map[augmentedRule, struct{}])
	
	for _, ar := range ars.GetAllKeys() {
		nextSymbol := ar.getNextSymbol()
		if nextSymbol == nil {
			continue
		}

		if _, ok := transitions[*nextSymbol]; !ok {
			transitions[*nextSymbol] = utils.NewMap[augmentedRule, struct{}]()
		}
		closure := getClosure(g, ar.shiftedCopy())
		transitions[*nextSymbol].PutAll(closure.GetAll())
	}
	return transitions
}

func findState(target *utils.Map[augmentedRule, struct{}], states []*lr1AutomationState) *lr1AutomationState {
	for _, state := range states {
		if target.KeysEqual(*state.augmentedRules) {
			return state
		}
	}
	return nil
}

func (g *grammar) generateLR1() (*lr1AutomationState, []*lr1AutomationState) {
	var id uint = 0

	kernel := newLR1AutomationState(&id)
	startRule := NewAugmentedRule(g.rules[0], 0, map[symbol]struct{}{EndOfInput: {}})
	kernel.augmentedRules = getClosure(g, startRule)

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