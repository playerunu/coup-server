package models

type Game struct {
	Players       []Player
	CurrentPlayer Player
	Deck          []Card
	PlayerActions []PlayerAction
}

func NewGame() *Game {
	game := &Game{
		Players:       []Player{},
		Deck:          newDeck(),
		PlayerActions: []PlayerAction{},
	}

	return game
}

func (game *Game) DrawCards(howMany int) []Card {
	var cards []Card
	for i := 0; i < howMany; i++ {
		cards = append(cards, game.Deck[len(game.Deck)-1])
	}

	game.Deck = game.Deck[:len(game.Deck)-howMany-1]

	return cards
}
