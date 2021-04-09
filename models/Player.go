package models

type Player struct {
	id    string `json:id`
	name  string `json:name`
	cards [2]Card
	coins int `json:coins`
}
