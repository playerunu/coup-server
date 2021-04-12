package models

import (
	"encoding/json"
	"math/rand"
	"time"
)

type Influence int

const (
	Duke Influence = iota
	Captain
	Assassin
	Contessa
	Ambassador
	length
)

type Card struct {
	influence  Influence
	isRevealed bool
}

func (card *Card) MarshalJSON() ([]byte, error) {
	marshalledCard := struct {
		Influence  Influence `json:"influence,omitempty"`
		IsRevealed bool      `json:"isRevealed"`
	}{
		IsRevealed: card.isRevealed,
	}

	if card.isRevealed {
		marshalledCard.Influence = card.influence
	}

	return json.Marshal(marshalledCard)
}

func newCard(influence Influence) *Card {
	return &Card{
		influence:  influence,
		isRevealed: false,
	}
}

func NewDeck() []Card {
	deck := []Card{}

	for influence := Influence(0); influence < Influence(length); influence++ {
		for i := 0; i < 3; i++ {
			deck = append(deck, *newCard(influence))
		}
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(deck), func(i, j int) { deck[i], deck[j] = deck[j], deck[i] })

	return deck
}
