package core

import (
	models "coup-server/model"
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type GameEngine struct {
	game                   *models.Game
	clientUpdatesChannel   chan ClientMessage
	clientsPrivateChannel  *chan ClientMessage
	globalBroadcastChannel *chan []byte
}

func NewGameEngine(globalBroadcastChannel *chan []byte, clientsPrivateChannel *chan ClientMessage) *GameEngine {
	return &GameEngine{
		game:                   models.NewGame(),
		clientUpdatesChannel:   make(chan ClientMessage),
		globalBroadcastChannel: globalBroadcastChannel,
		clientsPrivateChannel:  clientsPrivateChannel,
	}
}

func (engine *GameEngine) Run() {
	for {
		clientMessage := <-engine.clientUpdatesChannel
		rawMessage := *clientMessage.Payload

		var gameMessage GameMessage
		err := json.Unmarshal(rawMessage, &gameMessage)
		if err != nil {
			log.Fatalln("error:", err)
		}

		switch gameMessage.MessageType {
		case PlayerJoined:
			engine.onPlayerJoin(gameMessage, clientMessage.ClientUuid)
		}
	}
}

func (engine *GameEngine) ReadClientMessage(message ClientMessage) {
	engine.clientUpdatesChannel <- message
}

func (engine *GameEngine) SendClientMessage(player *models.Player, message []byte) {

}

func (engine *GameEngine) GlobalBroadcast(messageType MessageType) {
	gameJson, err := json.Marshal(engine.game)
	if err != nil {
		log.Fatal(err)
	}

	var gameMsg = GameMessage{
		MessageType: messageType,
		Data:        gameJson,
	}

	broadcastMessage, err := json.Marshal(gameMsg)
	if err != nil {
		log.Fatal(err)
	}

	*engine.globalBroadcastChannel <- broadcastMessage
}

func (engine *GameEngine) onPlayerJoin(message GameMessage, uuid uuid.UUID) {
	var player models.Player
	err := json.Unmarshal(message.Data, &player)
	if err != nil {
		log.Fatalln("error:", err)
	}
	player.SetConnectoinUuuid(uuid)

	engine.registerPlayer(player)
}

// Registers a new player
func (engine *GameEngine) registerPlayer(player models.Player) {
	// Draw 2 random cards
	copy(player.Cards[:], engine.game.DrawCards(2))
	// Give the initial coins
	engine.takeCoins(player.Name, 2)

	engine.game.Players = append(engine.game.Players, player)

	if len(engine.game.Players) >= 2 {
		engine.startGame()
	}
}

// Individually sends to each player its cards influences
func (engine *GameEngine) sendCardInfluences() {
	for playerIdx := range engine.game.Players {
		player := &engine.game.Players[playerIdx]

		fullCards := []models.MarshalledCard{}
		for cardIdx := range player.Cards {
			card := &player.Cards[cardIdx]
			fullCards = append(fullCards, card.MarshalCard(true))
		}

		fullCardsJson, err := json.Marshal(fullCards)
		if err != nil {
			log.Fatalln("error:", err)
		}

		var gameMsg = GameMessage{
			MessageType: YourCards,
			Data:        fullCardsJson,
		}

		gameMessageJson, err := json.Marshal(gameMsg)
		if err != nil {
			log.Fatal(err)
		}

		var clientMsg = ClientMessage{
			ClientUuid: player.GetConnectionUuid(),
			Payload:    &gameMessageJson,
		}

		*engine.clientsPrivateChannel <- clientMsg
	}
}

func (engine *GameEngine) startGame() {
	// Send each player info about its cards influences
	engine.sendCardInfluences()

	// Shuffle the players list
	rand.Seed(time.Now().UnixNano())
	players := engine.game.Players
	rand.Shuffle(len(players), func(i, j int) { players[i], players[j] = players[j], players[i] })

	// Assign each player its gamePosition
	// and set the first player to act
	for i := 0; i < len(players); i++ {
		players[i].GamePosition = i
	}
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
