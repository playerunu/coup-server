package models

type PlayerAction struct {
	Action           Action  `json:"action"`
	VsPlayer         *Player `json:"vsPlayer,omitempty"`
	ChallengedBy     *Player `json:"challengedBy,omitempty"`
	ChallengeSuccess *bool   `json:"challengeSuccess,omitempty"`
	BlockAction      *Block  `json:"blockAction,omitempty"`
}

func (playerAction *PlayerAction) IsSuccessful() bool {
	if !playerAction.Action.HasCounterAction {
		return true
	}

	if playerAction.ChallengedBy != nil {
		// The action was challenged, return the challenge result
		return *playerAction.ChallengeSuccess
	}

	if playerAction.BlockAction != nil {
		if playerAction.BlockAction.ChallengedBy == nil {
			// The action was blocked and nobody challenged it
			return false
		} else {
			// The block action was challenged, return
			return !*playerAction.BlockAction.ChallengeSuccess
		}
	}

	return true
}
