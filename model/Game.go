package models

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

// test
const (
	TOTAL_COINS              int = 50
	INITIAL_COINS_COUNT      int = 2
	ASSASSINATE_COINS_AMOUNT int = 3
	COUP_COINS_AMOUNT        int = 7
	STEAL_COINS_AMOUNT       int = 2
)

type Game struct {
	Players          []Player    `json:"players"`
	RemainingPlayers int         `json:"remainingPlayers"`
	Winner           *Player     `json:"winner:omitempty"`
	CurrentPlayer    *Player     `json:"currentPlayer"`
	CurrentMove      *PlayerMove `json:"currentMove,omitempty"`
	TableCoins       int         `json:"tableCoins"`
	deck             []Card
}

func NewGame() *Game {
	game := &Game{
		Players:    []Player{},
		deck:       NewDeck(),
		TableCoins: TOTAL_COINS,
	}

	game.shuffleDeck()

	return game
}

func (game *Game) DrawCards(howMany int) []Card {
	var cards []Card
	for i := 0; i < howMany; i++ {
		cards = append(cards, game.deck[len(game.deck)-1])
		game.deck = game.deck[:len(game.deck)-1]
	}

	return cards
}

func (game *Game) DrawCard() Card {
	return game.DrawCards(1)[0]
}

func (game *Game) InsertCard(card *Card) {
	game.deck = append(game.deck, *card)
	game.shuffleDeck()
}

func (game *Game) InsertCardAndDraw(card *Card) Card {
	game.InsertCard(card)

	return game.DrawCard()
}

func (game *Game) GetCoinsFromTable(player *Player, coinsAmount int) {
	if game.TableCoins < coinsAmount {
		log.Fatal("Not enough coins on the table")
	}

	player.Coins += coinsAmount
	game.TableCoins -= coinsAmount
}

func (game *Game) PutCoinsOnTable(player *Player, coinsAmount int) {
	player.Coins -= coinsAmount
	game.TableCoins += coinsAmount
}

func (game *Game) GetPlayerByName(playerName string) *Player {
	for index := range game.Players {
		player := &game.Players[index]
		if player.Name == playerName {
			return player
		}
	}

	return nil
}

func (game *Game) GetWinner() *Player {
	if game.RemainingPlayers > 1 {
		return nil
	}

	for index := range game.Players {
		player := &game.Players[index]
		if player.RemainingCards() > 0 {
			return player
		}
	}

	log.Fatal("Invalid game state! There is only one remaining player, but all player have 0 remaining cards")
	return nil
}

func (game *Game) ValidateState() bool {
	totalPlayerCoins := 0
	influencesCount := make(map[Influence]int)

	for playerIdx := range game.Players {
		player := game.Players[playerIdx]

		totalPlayerCoins += player.Coins

		influencesCount[player.Card1.GetInfluence()]++
		influencesCount[player.Card2.GetInfluence()]++
	}

	for cardIdx := range game.deck {
		card := game.deck[cardIdx]

		influencesCount[card.GetInfluence()]++
	}

	totalCoins := game.TableCoins + totalPlayerCoins
	if totalCoins != TOTAL_COINS {
		fmt.Println("Invalid total coins amount : ", totalCoins)
		return false
	}

	if len(influencesCount) != 5 {
		fmt.Println("Some influences are missing completly from game")
		return false
	}

	for influence, count := range influencesCount {
		if count != 3 {
			fmt.Println("There are only ", count, InfluenceToStr(influence), " cards")
			return false
		}
	}

	return true
}

func (game *Game) shuffleDeck() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(game.deck), func(i, j int) {
		game.deck[i], game.deck[j] = game.deck[j], game.deck[i]
	})
}
