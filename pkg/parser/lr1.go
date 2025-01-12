package parser

import (
	"fmt"
	"maps"
	"reflect"
	"strings"

	"asritha.dev/compiler/pkg/utils"
	."asritha.dev/compiler/pkg/grammar"
)

type arLookaheadMap map[augmentedRule]utils.Set[Symbol] //augmented rule -> lookahead

type lrAutomationState struct {
	id             uint
	arLookaheadMap
	transitions    map[Symbol]*lrAutomationState
}

func newLR1AutomationState(id *uint, arLookaheadMap arLookaheadMap) *lrAutomationState {
	*id++
	return &lrAutomationState{
		id:             *id - 1,
		arLookaheadMap: arLookaheadMap,
		transitions:    make(map[Symbol]*lrAutomationState),
	}
}

func (s *lrAutomationState) String() string {
	longestRule := 0
	for r := range s.arLookaheadMap {
		rLen := len(r.String())
		if rLen > longestRule {
			longestRule = rLen
		}
	}

	rules := make([]string, len(s.arLookaheadMap))
	i := 0
	for ar, lookahead := range s.arLookaheadMap {
		rules[i] = ar.StringWithLookahead(lookahead, longestRule+6)
		i++
	}
	return fmt.Sprintf("State %d\n%s", s.id, strings.Join(rules, "\n"))
}

func (ar augmentedRule) getClosureRecursion(g *Grammar, closure arLookaheadMap) {
	nextSymbol := ar.getNextSymbol()
	if nextSymbol == nil || nextSymbol.SymbolType != NonTermSymbol {
		return
	}

	sentential := ar.SententialForm[ar.position+1:]
	newLookahead := g.GenerateFirstSet(sentential...)
	if _, ok := newLookahead[Epsilon]; ok {
		utils.AddToMap(closure[ar], newLookahead) //add lookahead of current ar if it can be finished with epsilon
	}
	delete(newLookahead, Epsilon)

	for _, rule := range g.RuleNTMap[nextSymbol.Name] {
		newAR := *NewAugmentedRule(rule, 0)
		if _, ARExists := closure[newAR]; ARExists {
			utils.AddToMap(newLookahead, closure[newAR])
		} else {
			closure[newAR] = newLookahead
			newAR.getClosureRecursion(g, closure)
		}
	}
}

func getClosure(g *Grammar, initial arLookaheadMap) {
	for ar := range initial {
		ar.getClosureRecursion(g, initial)
	}
}

func getTransitions(g *Grammar, core arLookaheadMap) map[Symbol]arLookaheadMap {
	transitions := make(map[Symbol]arLookaheadMap)

	for ar, lookahead := range core {
		nextSymbol := ar.getNextSymbol()
		if nextSymbol == nil {
			continue
		}

		if _, ok := transitions[*nextSymbol]; !ok {
			transitions[*nextSymbol] = make(arLookaheadMap)
		}

		newAR := *NewAugmentedRule(ar.Rule, ar.position+1)
		//if transitions are not make correctly it may be HERE
		//removed check for if transitions[*nextSymbol][newAR] was already set because it seemed redundant
		transitions[*nextSymbol][newAR] = maps.Clone(lookahead)

		getClosure(g, transitions[*nextSymbol])
	}
	return transitions
}

func findState(target arLookaheadMap, states []*lrAutomationState) *lrAutomationState {
	for _, state := range states {
		if reflect.DeepEqual(target, state.arLookaheadMap) {
			return state
		}
	}
	return nil
}

func generateLR1(g *Grammar) (*lrAutomationState, []*lrAutomationState) {
	var id uint = 0

	kernel := newLR1AutomationState(&id, arLookaheadMap{
		*NewAugmentedRule(g.FirstRule, 0): {EndOfInput: struct{}{}},
	})
	getClosure(g, kernel.arLookaheadMap)

	states := []*lrAutomationState{kernel}
	openList := []*lrAutomationState{kernel}

	for len(openList) > 0 {
		state := openList[0]
		openList = openList[1:]

		for transition, rules := range getTransitions(g, state.arLookaheadMap) {
			if transitionState := findState(rules, states); transitionState != nil {
				state.transitions[transition] = transitionState
				continue
			}

			newState := newLR1AutomationState(&id, rules)
			state.transitions[transition] = newState
			states = append(states, newState)
			openList = append(openList, newState)
		}
	}

	return kernel, states
}

func findStateCore(target arLookaheadMap, states []*lrAutomationState) *lrAutomationState {
	for _, state := range states {
		if  utils.HasSameKeys(target, state.arLookaheadMap) {
			return state
		}
	}
	return nil
}

/*
Basically LR1 but merges states that have the same core (same augmented rules not including lookahead)
*/
func generateLALR(g *Grammar) (*lrAutomationState, []*lrAutomationState) {
	var id uint = 0

	kernel := newLR1AutomationState(&id, arLookaheadMap{
		*NewAugmentedRule(g.FirstRule, 0): {EndOfInput: struct{}{}},
	})
	getClosure(g, kernel.arLookaheadMap)

	states := []*lrAutomationState{kernel}
	openList := []*lrAutomationState{kernel}

	for len(openList) > 0 {
		state := openList[0]
		openList = openList[1:]

		for transition, rules := range getTransitions(g, state.arLookaheadMap) {
			if transitionState := findStateCore(rules, states); transitionState != nil {
				added := false
				for ar, lookahead := range transitionState.arLookaheadMap {
					if utils.AddToMap(rules[ar], lookahead) != 0 {
						added = true
					}
				}
				state.transitions[transition] = transitionState
				
				if added {
					openList = append(openList, transitionState)
				}
				
				continue
			}

			newState := newLR1AutomationState(&id, rules)
			state.transitions[transition] = newState
			states = append(states, newState)
			openList = append(openList, newState)
		}
	}

	return kernel, states
}

func makeMermaid(states []*lrAutomationState) string {
	mermaid := ""
	for _, state := range states {
		mermaid += fmt.Sprintf("state \"\n%s\n\" as s%d\n", strings.ReplaceAll(state.String(), "\"", "'"), state.id)
		for key, val := range state.transitions {
			mermaid += fmt.Sprintf("s%d-->s%d: '%s'\n", state.id, val.id, key.String())
		}
	}
	return mermaid
}

func makeGraphvizSafe(str string) string {
	return strings.ReplaceAll(
		strings.ReplaceAll(
			strings.ReplaceAll(
				str,
				`\`, `\\`,
			),
			"\n", `\l`,
		),
		`"`, `\"`,
	)
}

func makeGraphviz(states []*lrAutomationState) string {
	graph := "node [shape=box]\n"
	for _, state := range states {
		graph += fmt.Sprintf("s%d [label=\"%s\"]\n", state.id, makeGraphvizSafe(state.String()))
		for key, val := range state.transitions {
			graph += fmt.Sprintf("s%d -> s%d [label=\"%s\"]\n", state.id, val.id, makeGraphvizSafe(key.String()))
		}
	}
	return graph
}