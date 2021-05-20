package models

type ExchangeResult struct {
	Player         *Player   `json:"player"`
	NewPlayerCards *TwoCards `json:"playerCards"`
	DeckCards      *TwoCards `json:"deckCards"`
}
