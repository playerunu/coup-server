package core

import (
	models "coup-server/model"
	"log"
)

func DrawInfluence(engine *GameEngine, player *models.Player, influence models.Influence, cardType CardType) {
	var oldCard *models.Card

	switch cardType {
	case Card1:
		oldCard = &player.Card1
	case Card2:
		oldCard = &player.Card2
	case AnyUnrevealed:
		if !player.Card1.IsRevealed {
			oldCard = &player.Card1
			cardType = Card1
		} else if !player.Card2.IsRevealed {
			oldCard = &player.Card2
			cardType = Card2
		} else {
			log.Fatal("Both player's cards are revealed, can't draw a new card!")
		}
	}

	newCard := models.InsertAndDrawInfluence(engine.Game, influence, oldCard)
	if cardType == Card1 {
		player.Card1 = newCard
	} else if cardType == Card2 {
		player.Card2 = newCard
	}
}
