package models

import (
	"encoding/json"
)

type Card struct {
	influence  Influence
	IsRevealed bool
}

type MarshalledCard struct {
	Influence  Influence `json:"influence,omitempty"`
	IsRevealed bool      `json:"isRevealed"`
}

type TwoCards struct {
	Card1 MarshalledCard `json:"card1"`
	Card2 MarshalledCard `json:"card2"`
}

func (card *Card) MarshalCard(includeInfluence bool) MarshalledCard {
	marshalledCard := MarshalledCard{
		IsRevealed: card.IsRevealed,
	}

	if card.IsRevealed || includeInfluence {
		marshalledCard.Influence = card.influence
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

func (marshalledCard *MarshalledCard) ToCard() *Card {
	return &Card{
		influence: marshalledCard.Influence,
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

func newCard(influence Influence) *Card {
	return &Card{
		influence:  influence,
		IsRevealed: false,
	}
}
