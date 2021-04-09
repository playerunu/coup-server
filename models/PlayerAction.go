package models

type PlayerAction struct {
	action          Action
	vsPlayer        Player
	counterActrions []PlayerCounterAction
}
