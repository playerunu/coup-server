package main

type Influence int

const (
	Duke Influence = iota
	Captain
	Assassin
	Contessa
	Ambassador
)

type Card struct {
	influence  Influence
	isRevealed bool
}
