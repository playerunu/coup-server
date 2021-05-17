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
	ChallengeAction        = "CurrentActionChallenge"
	ChallenegeActionResult = "ChallengeActionResult"
	Block                  = "BlockAction"
	ChallengeBlock         = "ChallengeBlock"
	ChallengeBlockResult   = "ChallengeBlockResult"
	RevealCard             = "RevealCard"
	ActionResult           = "ActionResult"
	NextPlayer             = "NextPlayer"
	GameOver               = "GameOver"
)

type GameMessage struct {
	MessageType MessageType
	Data        json.RawMessage
}
