package models

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

func newCard(influence Influence) Card {
	card := Card{
		influence: influence,
		isRevealed: false,
	}
	
	return card
}

func newDeck() []Card {
	deck := []Card {}
	
	for influence := Influence(0); influence < Influence(length); influence++ {
		for i := 0; i < 3; i++ {
			deck = append(deck, newCard(influence));
		}
	}
	
	return deck;
}

