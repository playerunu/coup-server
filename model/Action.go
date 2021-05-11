package models

type ActionType int

const (
	TakeOneCoin ActionType = iota
	TakeTwoCoins
	TakeThreeCoins
	Exchange
	Assasinate
	Steal
	Coup
)

type Action struct {
	ActionType       ActionType `json:"actionType"`
	HasCounterAction bool       `json:"hasCounterAction"`
}

func NewAction(actionType ActionType) *Action {
	action := Action{
		ActionType:       actionType,
		HasCounterAction: true,
	}

	if actionType == TakeOneCoin {
		action.HasCounterAction = false
	}

	return &action
}
