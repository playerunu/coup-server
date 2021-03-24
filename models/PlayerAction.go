package main

type PlayerAction struct {
	action          Action
	vsPlayer        Player
	counterActrions []PlayerCounterAction
}
