package models

import "github.com/google/uuid"

type Player struct {
	Id             string  `json:"id"`
	Name           string  `json:"name"`
	Cards          [2]Card `json:"cards"`
	Coins          int     `json:"coins"`
	GamePosition   int     `json:"gamePostion"`
	connectionUuid uuid.UUID
}

func (player *Player) GetConnectionUuid() uuid.UUID {
	return player.connectionUuid
}

func (player *Player) SetConnectoinUuuid(connectionUuid uuid.UUID) {
	player.connectionUuid = connectionUuid
}
