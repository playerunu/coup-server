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
		case HeroPlayerAction:
			engine.onPlayerAction(gameMessage, clientMessage.ClientUuid)
		}

	}
}

func (engine *GameEngine) ReadClientMessage(message ClientMessage) {
	engine.clientUpdatesChannel <- message
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

func (engine *GameEngine) onPlayerAction(message GameMessage, uuid uuid.UUID) {
	var playerAction models.PlayerAction
	err := json.Unmarshal(message.Data, &playerAction)
	if err != nil {
		log.Fatalln("Error while unmarshalling game message: ", message.MessageType, err)
	}

	game := engine.game

	switch playerAction.Action.ActionType {
	case models.TakeOneCoin:
		engine.takeCoins(game.CurrentPlayer.Name, 1)
	case models.TakeTwoCoins:
		engine.takeCoins(game.CurrentPlayer.Name, 1)
	case models.TakeThreeCoins:
		engine.takeCoins(game.CurrentPlayer.Name, 3)
	}

	engine.game.CurrentPlayerAction = &playerAction
	engine.GlobalBroadcast(PlayerAction)
}

// Registers a new player
func (engine *GameEngine) registerPlayer(player models.Player) {
	// Draw 2 random cards
	player.Card1 = engine.game.DrawCard()
	player.Card2 = engine.game.DrawCard()

	engine.game.Players = append(engine.game.Players, player)

	// Give the initial coins
	engine.takeCoins(player.Name, 2)

	if len(engine.game.Players) >= 2 {
		engine.startGame()
	}
}

// Individually sends to each player its cards influences
func (engine *GameEngine) sendCardInfluences() {
	for playerIdx := range engine.game.Players {
		player := &engine.game.Players[playerIdx]

		fullCard1 := player.Card1.MarshalCard(true)
		fullCard2 := player.Card2.MarshalCard(true)

		fullCardsJson, err := json.Marshal(struct {
			Card1 models.MarshalledCard `json:"card1"`
			Card2 models.MarshalledCard `json:"card2"`
		}{
			Card1: fullCard1,
			Card2: fullCard2,
		})
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
			return
		}
	}
}
