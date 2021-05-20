package models

type Challenge struct {
	ChallengedBy  *Player `json:"challengedBy,omitempty"`
	Success       *bool   `json:"challengeStatus,omitempty"`
	WaitingReveal bool    `json:"waitingReveal,omitempty"`
}

func NewChallenge(challengedBy *Player) *Challenge {
	return &Challenge{
		ChallengedBy:  challengedBy,
		WaitingReveal: true,
	}
}
