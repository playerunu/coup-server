package models

import "github.com/google/uuid"

type Player struct {
	Name           string `json:"name"`
	Card1          Card   `json:"card1"`
	Card2          Card   `json:"card2"`
	Coins          int    `json:"coins"`
	GamePosition   int    `json:"gamePosition"`
	connectionUuid uuid.UUID
}

func (player *Player) GetConnectionUuid() uuid.UUID {
	return player.connectionUuid
}

func (player *Player) SetConnectionUuid(connectionUuid uuid.UUID) {
	player.connectionUuid = connectionUuid
}

func (player *Player) HasInfluence(influence Influence) bool {
	return (!player.Card1.IsRevealed() && player.Card1.GetInfluence() == influence) ||
		(!player.Card2.IsRevealed() && player.Card2.GetInfluence() == influence)
}
