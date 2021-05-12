package models

type PlayerMove struct {
	Action          Action     `json:"action"`
	WaitingReveal   *bool      `json:"waitingReveal,omitempty"`
	WaitingExchange *bool      `json:"waitinExchange,omitempty"`
	VsPlayer        *Player    `json:"vsPlayer,omitempty"`
	Challenge       *Challenge `json:"challenge,omitempty"`
	Block           *Block     `json:"blockAction,omitempty"`
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
		if playerMove.Block.Challenge.ChallengedBy == nil {
			// The action was blocked and nobody challenged it
			return false
		} else {
			// The block action was challenged, check the block action result
			return !*playerMove.Block.Challenge.Success
		}
	}

	return true
}

// Checks if the current player action or its block
// can still be countered by a block or a challenge
func (playerMove *PlayerMove) CanCounter() bool {
	if !playerMove.Action.CanCounter() {
		return false
	}

	if playerMove.Action.CanBlock && playerMove.Block == nil {
		return true
	}

	if playerMove.Action.CanChallenge && playerMove.Challenge == nil &&
		// An action that has been blocked cannot be challenged anymore
		playerMove.Block == nil {
		return true
	}

	if playerMove.Block != nil && playerMove.Block.Challenge == nil {
		return true
	}

	// This should never be reached
	return false
}
