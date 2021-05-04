package models

type Block struct {
	Player              Player    `json:"player"`
	PretendingInfluence Influence `json:"pretendingInfluence"`
	ChallengedBy        Player    `json:"challengedBy"`
}
