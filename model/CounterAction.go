package models

type CounterActionType int

const (
	Block CounterActionType = iota
	Challenge
)

type CounterAction struct {
	CounterActionType CounterActionType `json:"counterActionType"`
	HasCounterAction  bool              `json:"hasCounterAction"`
}
