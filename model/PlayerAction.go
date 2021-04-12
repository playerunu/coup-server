package models

type PlayerAction struct {
	Action         Action                `json:"action"`
	VsPlayer       Player                `json:"vsPlayer"`
	CounterActions []PlayerCounterAction `json:"counterActions"`
}
