package models

import "github.com/google/uuid"

type Player struct {
	Id             string `json:"id"`
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

func (player *Player) SetConnectoinUuuid(connectionUuid uuid.UUID) {
	player.connectionUuid = connectionUuid
}
