package models

type PlayerCounterAction struct {
	CounterActionType   CounterActionType `json:"counterActionType"`
	Player              Player            `json:"player"`
	VsPlayer            Player            `json:"vsPlayer"`
	PretendingInfluence Influence         `json:"pretendingInfluence"`
}
