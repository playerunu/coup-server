package models

type PlayerAction struct {
	Action       Action `json:"action"`
	VsPlayer     Player `json:"vsPlayer"`
	ChallengedBy Player `json:"challengedBy"`
	BlockAction  Block  `json:"blockAction"`
}
