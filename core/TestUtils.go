package core

import (
	models "coup-server/model"
	"log"
)

func DrawInfluence(engine *GameEngine, player *models.Player, influence models.Influence, cardType models.CardType) {
	var oldCard *models.Card

	switch cardType {
	case models.Card1:
		oldCard = &player.Card1
	case models.Card2:
		oldCard = &player.Card2
	case models.AnyUnrevealed:
		if !player.Card1.IsRevealed {
			oldCard = &player.Card1
			cardType = models.Card1
		} else if !player.Card2.IsRevealed {
			oldCard = &player.Card2
			cardType = models.Card2
		} else {
			log.Fatal("Both player's cards are revealed, can't draw a new card!")
		}
	}

	newCard := models.InsertAndDrawInfluence(engine.Game, influence, oldCard)
	if cardType == models.Card1 {
		player.Card1 = newCard
	} else if cardType == models.Card2 {
		player.Card2 = newCard
	}
}
