package models

type Game struct {
	players       []Player `json:"players"`
	currentPlayer Player `json:currentPlayer`
	deck          []Card 
	playerActions []PlayerAction 
}

func NewGame() Game {
	game := Game{
		players : []Player{},
		deck: newDeck(),
		playerActions: []PlayerAction{},
	}
	
	return game
}

