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
	ActionType   ActionType `json:"actionType"`
	CanChallenge bool       `json:"canChallenge"`
	CanBlock     bool       `json:"canBlock"`
	influence    *Influence
}

func NewAction(actionType ActionType) *Action {
	action := Action{
		ActionType:   actionType,
		CanChallenge: true,
		CanBlock:     true,
		influence:    nil,
	}

	if actionType == TakeOneCoin || actionType == Coup || actionType == TakeTwoCoins {
		action.CanChallenge = false
	}

	if actionType == TakeOneCoin || actionType == Coup || actionType == Exchange {
		action.CanBlock = false
	}

	var influence Influence
	switch actionType {
	case TakeThreeCoins:
		influence = Duke
	case Exchange:
		influence = Ambassador
	case Steal:
		influence = Captain
	case Assasinate:
		influence = Assassin
	}

	action.influence = &influence

	return &action
}

func (action *Action) CanCounter() bool {
	return action.CanBlock || action.CanChallenge
}

func (action *Action) GetInfluence() *Influence {
	return action.influence
}
