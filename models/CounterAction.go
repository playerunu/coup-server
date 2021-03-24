package main

type CounterActionType int

const (
	Block CounterActionType = iota
	Challenge
)

type CounterAction struct {
	counterActionType CounterActionType
	hasCounterAction  bool
}
