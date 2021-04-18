package core

import (
	"encoding/json"
)

type MessageType string

const (
	// Client messages
	HeroPlayerAction = "HeroPlayerAction"

	// Server messages
	PlayerJoined = "PlayerJoined"
	GameStarted  = "GameStarted"
	YourCards    = "YourCards"
	PlayerAction = "PlayerAction"
)

type GameMessage struct {
	MessageType MessageType
	Data        json.RawMessage
}
