package models

type PlayerAction struct {
	Action           Action  `json:"action"`
	VsPlayer         *Player `json:"vsPlayer,omitempty"`
	ChallengedBy     *Player `json:"challengedBy,omitempty"`
	ChallengeSuccess *bool   `json:"challengeSuccess,omitempty"`
	BlockAction      *Block  `json:"blockAction,omitempty"`
}
