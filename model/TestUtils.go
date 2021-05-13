package models

func InsertAndDrawInfluence(game *Game, influnece Influence, card *Card) Card {
	var newCard Card
	var newCardIndex int

	game.deck = append(game.deck, *card)
	for index := range game.deck {
		newCard = game.deck[index]
		newCardIndex = index
		if newCard.GetInfluence() == influnece {
			break
		}
	}

	game.deck = append(game.deck[:newCardIndex], game.deck[newCardIndex+1:]...)

	return newCard
}
