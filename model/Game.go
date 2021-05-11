package models

import (
	"math/rand"
	"time"
)

var TOTAL_COINS int = 50

type Game struct {
	Players             []Player      `json:"players"`
	CurrentPlayer       *Player       `json:"currentPlayer"`
	CurrentPlayerAction *PlayerAction `json:"currentPlayerAction,omitempty"`
	TableCoins          int           `json:"tableCoins"`
	deck                []Card
}

func NewGame() *Game {
	game := &Game{
		Players:             []Player{},
		deck:                NewDeck(),
		TableCoins:          TOTAL_COINS,
		CurrentPlayerAction: nil,
	}

	game.shuffleDeck()

	return game
}

func (game *Game) DrawCards(howMany int) []Card {
	var cards []Card
	for i := 0; i < howMany; i++ {
		cards = append(cards, game.deck[len(game.deck)-1])
		game.deck = game.deck[:len(game.deck)-1]
	}

	return cards
}

func (game *Game) DrawCard() Card {
	return game.DrawCards(1)[0]
}

func (game *Game) InsertAndDraw(card Card) Card {
	game.deck = append(game.deck, card)
	game.shuffleDeck()

	return game.DrawCard()

}

func (game *Game) shuffleDeck() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(game.deck), func(i, j int) {
		game.deck[i], game.deck[j] = game.deck[j], game.deck[i]
	})

}
