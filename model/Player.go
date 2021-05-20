package models

import (
	"log"

	"github.com/google/uuid"
)

type Player struct {
	Name           string `json:"name"`
	Card1          Card   `json:"card1"`
	Card2          Card   `json:"card2"`
	Coins          int    `json:"coins"`
	GamePosition   int    `json:"gamePosition"`
	connectionUuid uuid.UUID
}

type CardType int

const (
	Card1 CardType = iota
	Card2
	LastUnrevealed
	AnyUnrevealed
)

func (player *Player) GetConnectionUuid() uuid.UUID {
	return player.connectionUuid
}

func (player *Player) SetConnectionUuid(connectionUuid uuid.UUID) {
	player.connectionUuid = connectionUuid
}

func (player *Player) HasInfluence(influence Influence) bool {
	return (!player.Card1.IsRevealed && player.Card1.GetInfluence() == influence) ||
		(!player.Card2.IsRevealed && player.Card2.GetInfluence() == influence)
}

func (player *Player) RemainingCards() int {
	remaining := 2

	if player.Card1.IsRevealed {
		remaining -= 1
	}
	if player.Card2.IsRevealed {
		remaining -= 1
	}

	return remaining
}

func (player *Player) RevealCard(cardType CardType) {
	if cardType == AnyUnrevealed {
		if !player.Card1.IsRevealed {
			cardType = Card1
		} else if !player.Card2.IsRevealed {
			cardType = Card2
		} else {
			log.Fatal("Both player cards are revealed and we are supposed to find one unrevealed card")
		}
	}

	switch cardType {
	case Card1:
		player.Card1.Reveal()
	case Card2:
		player.Card2.Reveal()
	case LastUnrevealed:
		player.RevealLastCard()
	}
}

func (player *Player) RevealLastCard() {
	if player.RemainingCards() != 1 {
		return
	}

	if !player.Card1.IsRevealed {
		player.Card1.Reveal()
	}
	if !player.Card2.IsRevealed {
		player.Card2.Reveal()
	}
}

func (player *Player) IsEliminated() bool {
	return player.RemainingCards() == 0
}

func (player *Player) StealFromPlayer(vsPlayer *Player) {
	if vsPlayer.Coins >= STEAL_COINS_AMOUNT {
		player.Coins += STEAL_COINS_AMOUNT
		vsPlayer.Coins -= STEAL_COINS_AMOUNT
	} else {
		player.Coins += vsPlayer.Coins
		vsPlayer.Coins = 0
	}
}
