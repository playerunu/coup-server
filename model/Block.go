package models

type Block struct {
	Player              *Player    `json:"player"`
	PretendingInfluence *Influence `json:"pretendingInfluence"`
	Challenge           *Challenge `json:"challenge,omitempty"`
}
