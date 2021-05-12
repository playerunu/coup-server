package models

type Challenge struct {
	ChallengedBy  *Player `json:"challengedBy,omitempty"`
	Success       *bool   `json:"challengeStatus,omitempty"`
	WaitingReveal bool    `json:"waitingReveal"`
}
