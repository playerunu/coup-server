package main

type Game struct {
	players       []Player
	currentPlayer Player
	deck          []Card
	playerActions []PlayerAction
}
