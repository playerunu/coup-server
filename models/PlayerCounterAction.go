package main

type PlayerCounterAction struct {
	counterActionType   CounterActionType
	player              Player
	vsPlayer            Player
	pretendingInfluence Influence
}
