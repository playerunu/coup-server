package core

import (
	"encoding/json"
)

type MessageType string

const (
	PlayerJoined          = "PlayerJoined"
	GameStarted           = "GameStarted"
	YourCards             = "YourCards"
	Action                = "Action"
	ActionResult          = "ActionResult"
	ChallengeAction       = "CurrentActionChallenge"
	ChallengeActionResult = "ChallengeActionResult"
	Block                 = "BlockAction"
	ChallengeBlock        = "ChallengeBlock"
	ChallengeBlockResult  = "ChallengeBlockResult"
	WaitingReveal         = "WaitingReveal"
	RevealCard            = "RevealCard"
	WaitingExchange       = "WaitingExchange"
	YourExchangeCards     = "YourExchangeCards"
	ExchangeComplete      = "ExchangeComplete"
	NextPlayer            = "NextPlayer"
	GameOver              = "GameOver"
)

type GameMessage struct {
	MessageType MessageType
	Data        json.RawMessage
}
