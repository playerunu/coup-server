package core

import (
	models "coup-server/model"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"
)

type GameEngine struct {
	game                   *models.Game
	clientUpdatesChannel   chan []byte
	globalBroadcastChannel *chan []byte
}

func NewGameEngine(globalBroadcastChannel *chan []byte) *GameEngine {
	return &GameEngine{
		game:                   models.NewGame(),
		clientUpdatesChannel:   make(chan []byte),
		globalBroadcastChannel: globalBroadcastChannel,
	}
}

func (engine *GameEngine) Run() {
	for {
		rawMessage := <-engine.clientUpdatesChannel

		var gameMessage GameMessage
		err := json.Unmarshal(rawMessage, &gameMessage)
		if err != nil {
			log.Fatalln("error:", err)
		}

		fmt.Println(gameMessage.MessageType)

		switch gameMessage.MessageType {
		case PlayerJoined:
			engine.onPlayerJoin(gameMessage)
		}
	}
}

func (engine *GameEngine) OnClientMessage(message []byte) {
	engine.clientUpdatesChannel <- message
}

func (engine *GameEngine) GlobalBroadcast(messageType MessageType) {
	gameJson, err := json.Marshal(engine.game)
	if err != nil {
		log.Fatal(err)
	}

	var gameStartMsg = GameMessage{
		MessageType: messageType,
		Data:        gameJson,
	}

	broadcastMessage, err := json.Marshal(gameStartMsg)
	if err != nil {
		log.Fatal(err)
	}

	*engine.globalBroadcastChannel <- broadcastMessage
}

func (engine *GameEngine) onPlayerJoin(message GameMessage) {
	var player models.Player
	err := json.Unmarshal(message.Data, &player)
	if err != nil {
		log.Fatalln("error:", err)
	}

	engine.registerPlayer(player)
}

// Registers a new player
func (engine *GameEngine) registerPlayer(player models.Player) {

	engine.game.Players = append(engine.game.Players, player)
	// Draw 2 random cards
	copy(player.Cards[:], engine.game.DrawCards(2))
	// Give the initial coins
	engine.takeCoins(player.Name, 2)

	if len(engine.game.Players) >= 2 {
		engine.startGame()
	}
}

func (engine *GameEngine) startGame() {
	// Shuffle the players list
	rand.Seed(time.Now().UnixNano())
	players := engine.game.Players
	rand.Shuffle(len(players), func(i, j int) { players[i], players[j] = players[j], players[i] })
	// Assign each player its gamePosition
	for i := 0; i < len(players); i++ {
		players[i].GamePosition = i
	}
	// Set the first player to act
	engine.game.CurrentPlayer = players[0]

	engine.GlobalBroadcast(GameStarted)
}

func (engine *GameEngine) takeCoins(playerName string, coinsAmount int) {
	if engine.game.TableCoins < coinsAmount {
		log.Fatal("Not enough coins on the table")
	}

	for i := 0; i < len(engine.game.Players); i++ {
		player := &engine.game.Players[i]
		if player.Name == playerName {
			player.Coins += coinsAmount
			engine.game.TableCoins -= coinsAmount
		}
	}
}
