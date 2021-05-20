package core

import (
	"encoding/json"
)

type MessageType string

const (
	PlayerJoined           = "PlayerJoined"
	GameStarted            = "GameStarted"
	YourCards              = "YourCards"
	Action                 = "Action"
	ActionResult           = "ActionResult"
	ChallengeAction        = "CurrentActionChallenge"
	ChallenegeActionResult = "ChallengeActionResult"
	Block                  = "BlockAction"
	ChallengeBlock         = "ChallengeBlock"
	ChallengeBlockResult   = "ChallengeBlockResult"
	RevealCard             = "RevealCard"
	ExchangeDeckCards      = "ExchangeDeckCards"
	ExchangeComplete       = "ExchangeComplete"
	NextPlayer             = "NextPlayer"
	GameOver               = "GameOver"
)

type GameMessage struct {
	MessageType MessageType
	Data        json.RawMessage
}
