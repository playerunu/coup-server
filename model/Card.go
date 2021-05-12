package models

import (
	"encoding/json"
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
	IsRevealed bool
}

type MarshalledCard struct {
	Influence  string `json:"influence,omitempty"`
	IsRevealed bool   `json:"isRevealed"`
}

func (card *Card) MarshalCard(includeInfluence bool) MarshalledCard {
	marshalledCard := MarshalledCard{
		IsRevealed: card.IsRevealed,
	}

	if card.IsRevealed || includeInfluence {
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

func (card *Card) GetInfluence() Influence {
	return card.influence
}

func (card *Card) Reveal() {
	card.IsRevealed = true
}

func newCard(influence Influence) *Card {
	return &Card{
		influence:  influence,
		IsRevealed: false,
	}
}

func NewDeck() []Card {
	deck := []Card{}

	for influence := Influence(0); influence < Influence(length); influence++ {
		for i := 0; i < 3; i++ {
			deck = append(deck, *newCard(influence))
		}
	}

	return deck
}
