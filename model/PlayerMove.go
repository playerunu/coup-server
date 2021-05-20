package models

import "log"

type PlayerMove struct {
	Action          Action     `json:"action"`
	Finished        bool       `json:"finished"`
	WaitingReveal   bool       `json:"waitingReveal"`
	WaitingExchange bool       `json:"waitinExchange"`
	VsPlayer        *Player    `json:"vsPlayer,omitempty"`
	Challenge       *Challenge `json:"challenge,omitempty"`
	Block           *Block     `json:"blockAction,omitempty"`
}

func NewPlayerMove(actionType ActionType, vsPlayer *Player) *PlayerMove {
	return &PlayerMove{
		Action:   *NewAction(actionType),
		Finished: false,
		VsPlayer: vsPlayer,
	}
}

func (playerMove *PlayerMove) IsSuccessful() bool {
	if !playerMove.Action.CanCounter() {
		return true
	}

	if playerMove.Challenge != nil {
		// The action was challenged, return the challenge result
		return *playerMove.Challenge.Success
	}

	if playerMove.Block != nil {
		if playerMove.Block.Challenge == nil {
			// The action was blocked and nobody challenged it
			return false
		} else {
			// The block action was challenged, check the challenge block action result
			return *playerMove.Block.Challenge.Success
		}
	}

	return true
}

// Checks if the current player action or its block
// can still be countered by a block or a challenge
func (playerMove *PlayerMove) CanCounter() bool {
	// The action has no counters or the move is already finished
	if !playerMove.Action.CanCounter() || playerMove.Finished {
		return false
	}

	// The action can be blocked
	if playerMove.Action.CanBlock {
		return playerMove.Block == nil
	}

	// The action ca be challlenged
	// Note that an action that has been blocked cannot be challenged anymore
	if playerMove.Action.CanChallenge && playerMove.Block == nil {
		return playerMove.Challenge == nil
	}

	// The block action can be challenged
	if playerMove.Block != nil && playerMove.Block.Challenge == nil {
		return playerMove.Block.Challenge == nil
	}

	// This should never be reached
	log.Fatal("Invalid game state while checking current move counters!")
	return false
}

func (playerMove *PlayerMove) IsWaitingMoveReveal() bool {
	return playerMove.WaitingReveal
}

func (playerMove *PlayerMove) IsWaitingChallengeReveal() bool {
	return playerMove.Challenge != nil && playerMove.Challenge.WaitingReveal
}

func (playerMove *PlayerMove) IsWaitingBlockReveal() bool {
	return playerMove.Block != nil && playerMove.Block.IsWaitingReveal()
}
