package models

type Player struct {
	Id    string
	Name  string
	Cards [2]Card
	Coins int
}
