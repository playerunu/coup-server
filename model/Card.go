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

type MarshalledCard struct {
	Influence  string `json:"influence,omitempty"`
	IsRevealed bool   `json:"isRevealed"`
}

func (card *Card) MarshalCard(includeInfluence bool) MarshalledCard {
	marshalledCard := MarshalledCard{
		IsRevealed: card.isRevealed,
	}

	if card.isRevealed || includeInfluence {
		switch card.influence {
		case Duke:
			marshalledCard.Influence = "Duke"
		case Captain:
			marshalledCard.Influence = "Captain"
		case Assassin:
			marshalledCard.Influence = "Assassin"
		case Contessa:
			marshalledCard.Influence = "Contessa"
		case Ambassador:
			marshalledCard.Influence = "Ambassador"
		}

	}

	return marshalledCard
}

func (card *Card) MarshalJSON() ([]byte, error) {
	return json.Marshal(card.MarshalCard(false))
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
