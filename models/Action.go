package main

type ActionType int

const (
	TakeOneCoin ActionType = iota
	TakeTwoCoins
	TakeThreeCoins
	Exchange
	Assasinate
	Steal
)

type Action struct {
	actionType       ActionType
	hasCounterAction bool
}

func newAction(actionType ActionType) *Action {
	action := Action{actionType: actionType, hasCounterAction: true}
	if actionType == TakeOneCoin {
		action.hasCounterAction = false
	}

	return &action
}
