package models

var TOTAL_COINS int = 50

type Game struct {
	Players       []Player       `json:"players"`
	CurrentPlayer Player         `json:"currentPlayer"`
	TableCoins    int            `json:"tableCoins"`
	PlayerActions []PlayerAction `json:"playerActions"`
	deck          []Card
}

func NewGame() *Game {
	game := &Game{
		Players:       []Player{},
		deck:          NewDeck(),
		PlayerActions: []PlayerAction{},
		TableCoins:    TOTAL_COINS,
	}

	return game
}

func (game *Game) DrawCards(howMany int) []Card {
	var cards []Card
	for i := 0; i < howMany; i++ {
		cards = append(cards, game.deck[len(game.deck)-1])
	}

	game.deck = game.deck[:len(game.deck)-howMany-1]

	return cards
}
