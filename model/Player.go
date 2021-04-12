package models

type Player struct {
	Id           string  `json:"id"`
	Name         string  `json:"name"`
	Cards        [2]Card `json:"cards"`
	Coins        int     `json:"coins"`
	GamePosition int     `json:"gamePostion"`
}
