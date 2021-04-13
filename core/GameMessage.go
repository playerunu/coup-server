package core

import (
	"encoding/json"
)

type MessageType string

const (
	PlayerJoined = "PlayerJoined"
	YourCards    = "YourCards"
	GameStarted  = "GameStarted"
)

type GameMessage struct {
	MessageType MessageType
	Data        json.RawMessage
}
