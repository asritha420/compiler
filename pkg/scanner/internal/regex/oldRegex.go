package regex

import (
	"fmt"
)

func convertRegexToNfaRecursion(regexASTRootNode RExpr, idToState map[uint]*fsm.FAState, id *uint) (*fsm.FAState, *fsm.FAState, error) {
	switch rootNode := regexASTRootNode.(type) {
	case *Concatenation:
		lNFAState, lNFALastState, err := convertRegexToNfaRecursion(rootNode.Left, idToState, id)
		if err != nil {
			return nil, nil, fmt.Errorf("concatination left node: %w", err)
		}
		rNFAState, rNFALastState, err := convertRegexToNfaRecursion(rootNode.Right, idToState, id)
		if err != nil {
			return nil, nil, fmt.Errorf("concatenation right node: %w", err)
		}

		lNFALastState.AddTransition(fsm.Epsilon, rNFAState)

		return lNFAState, rNFALastState, nil
	case *Alternation:
		lNFAState, lNFALastState, err := convertRegexToNfaRecursion(rootNode.Left, idToState, id)
		if err != nil {
			return nil, nil, fmt.Errorf("alternation left node: %w", err)
		}
		rNFAState, rNFALastState, err := convertRegexToNfaRecursion(rootNode.Right, idToState, id)
		if err != nil {
			return nil, nil, fmt.Errorf("alternation right node: %w", err)
		}

		start := fsm.NewNFAState(id, false)
		idToState[*id] = start

		end := fsm.NewNFAState(id, false)
		idToState[*id] = end

		start.AddTransition(fsm.Epsilon, lNFAState, rNFAState)
		lNFALastState.AddTransition(fsm.Epsilon, end)
		rNFALastState.AddTransition(fsm.Epsilon, end)

		return start, end, nil
	case *KleeneStar:
		NFAStartState, NFALastState, err := convertRegexToNfaRecursion(rootNode.Left, idToState, id)
		if err != nil {
			return nil, nil, fmt.Errorf("kleene star child node: %w", err)
		}

		start := fsm.NewNFAState(id, false)
		idToState[*id] = start

		end := fsm.NewNFAState(id, false)
		idToState[*id] = end

		start.AddTransition(fsm.Epsilon, NFAStartState, end)
		NFALastState.AddTransition(fsm.Epsilon, end)
		end.AddTransition(fsm.Epsilon, start)

		return start, end, nil
	case *Const:
		start := fsm.NewNFAState(id, false)
		idToState[*id] = start

		end := fsm.NewNFAState(id, false)
		idToState[*id] = end

		start.AddTransition(rootNode.Value, end)

		return start, end, nil
	default:
		return nil, nil, fmt.Errorf("invalid node: %v", rootNode)
	}
}

func ConvertRegexToNfa(regexASTRootNode RExpr) (*fsm.FAState, *fsm.FAState, map[uint]*fsm.FAState, error) {
	idMap := make(map[uint]*fsm.FAState)
	var idVar uint
	start, end, err := convertRegexToNfaRecursion(regexASTRootNode, idMap, &idVar)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("convertRegexToNfa: Unable to convert regex AST to NFA. Node trace:\n\t%w", err)
	}
	end.SetAccepting(true)
	return start, end, idMap, nil
}
