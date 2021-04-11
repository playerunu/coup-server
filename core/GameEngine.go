package core

import (
	"coup-server/models"
	"encoding/json"
	"fmt"
	"log"
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

func (engine *GameEngine) SendMessage(message []byte) {
	engine.clientUpdatesChannel <- message
}

func (engine *GameEngine) onPlayerJoin(message GameMessage) {
	var player models.Player
	err := json.Unmarshal(message.Data, &player)
	if err != nil {
		log.Fatalln("error:", err)
	}
	fmt.Println(player.Name)

	// Draw 2 random cards
	copy(player.Cards[:], engine.game.DrawCards(2))
	// Give the initial coins
	player.Coins = 2

	engine.game.Players = append(engine.game.Players, player)

	if len(engine.game.Players) >= 2 {
		var gameStartMsg = GameMessage{MessageType: GameStarted}
		b, err := json.Marshal(gameStartMsg)
		if err != nil {
			log.Fatal(err)
		}

		*engine.globalBroadcastChannel <- b
	}
}
