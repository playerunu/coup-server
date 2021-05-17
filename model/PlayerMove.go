package models

type PlayerMove struct {
	Action          Action     `json:"action"`
	Finished        bool       `json:"finished"`
	WaitingReveal   *bool      `json:"waitingReveal,omitempty"`
	WaitingExchange *bool      `json:"waitinExchange,omitempty"`
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
	if !playerMove.Action.CanCounter() {
		return false
	}

	if playerMove.Finished {
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

func (playerMove *PlayerMove) IsWaitingMoveReveal() bool {
	return playerMove.WaitingReveal != nil && *playerMove.WaitingReveal
}

func (playerMove *PlayerMove) IsWaitingChallengeReveal() bool {
	return playerMove.Challenge != nil && playerMove.Challenge.IsWaitingReveal()
}

func (playerMove *PlayerMove) IsWaitingBlockReveal() bool {
	return playerMove.Block != nil && playerMove.Block.IsWaitingReveal()
}
