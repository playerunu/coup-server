package core

import (
	models "coup-server/model"
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

const MAX_PLAYERS = 2
const INITIAL_COINS_COUNT = 2

type GameEngine struct {
	Game                   *models.Game
	waitingCountersTimer   *time.Timer
	clientUpdatesChannel   chan ClientMessage
	clientsPrivateChannel  *chan ClientMessage
	globalBroadcastChannel *chan []byte
}

func NewGameEngine(globalBroadcastChannel *chan []byte, clientsPrivateChannel *chan ClientMessage) *GameEngine {
	return &GameEngine{
		Game:                   models.NewGame(),
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
		case Action:
			engine.onPlayerAction(gameMessage, clientMessage.ClientUuid)
		case ChallengeAction:
		case Block:
		case ChallengeBlock:
			engine.onPlayerCounter(gameMessage, clientMessage.ClientUuid)
		case RevealCard:
			//engine.onCardReveal(gameMessage, clientMessage.ClientUuid)
		}
	}
}

func (engine *GameEngine) ReadClientMessage(message ClientMessage) {
	engine.clientUpdatesChannel <- message
}

func (engine *GameEngine) GlobalBroadcast(messageType MessageType) {
	gameJson, err := json.Marshal(engine.Game)
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

	game := engine.Game

	switch playerAction.Action.ActionType {
	case models.TakeOneCoin:
		engine.takeCoins(game.CurrentPlayer.Name, 1)
	case models.TakeTwoCoins:
		engine.takeCoins(game.CurrentPlayer.Name, 2)
	case models.TakeThreeCoins:
		engine.takeCoins(game.CurrentPlayer.Name, 3)
	case models.Assasinate:
	case models.Steal:
	case models.Exchange:
		// Those actions will only be broadcasted
		break
	}

	engine.Game.CurrentPlayerAction = &playerAction
	engine.GlobalBroadcast(Action)

	if playerAction.Action.HasCounterAction {
		engine.waitForCounters()
	} else {
		engine.nextPlayer()
	}
}

func (engine *GameEngine) onPlayerCounter(message GameMessage, uuid uuid.UUID) {
	// Always make sure we first stop the waiting counters timer
	engine.waitingCountersTimer.Stop()

	var playerAction models.PlayerAction
	err := json.Unmarshal(message.Data, &playerAction)
	if err != nil {
		log.Fatalln("Error while unmarshalling game message: ", message.MessageType, err)
	}

	switch message.MessageType {
	case ChallengeAction:
	case ChallengeBlock:
		engine.solveChallenge(message.MessageType)
	case Block:
		// Block counter only gets broadcasted
		break
	}

	engine.Game.CurrentPlayerAction = &playerAction
	engine.GlobalBroadcast(message.MessageType)

	// Restart the waiting counters timer only when blocking
	// For challenges, we restart the counter after the card is revealed
	if message.MessageType == Block {
		engine.waitForCounters()
	}
}

// Registers a new player
func (engine *GameEngine) registerPlayer(player models.Player) {
	// Draw 2 random cards
	player.Card1 = engine.Game.DrawCard()
	player.Card2 = engine.Game.DrawCard()

	engine.Game.Players = append(engine.Game.Players, player)

	// Give the initial coins
	engine.takeCoins(player.Name, INITIAL_COINS_COUNT)

	if len(engine.Game.Players) >= MAX_PLAYERS {
		engine.startGame()
	}
}

// Individually sends to each player its cards influences
func (engine *GameEngine) sendCardInfluences() {
	for playerIdx := range engine.Game.Players {
		engine.sendPlayerCardInfluences(&engine.Game.Players[playerIdx])
	}
}

func (engine *GameEngine) sendPlayerCardInfluences(player *models.Player) {
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

func (engine *GameEngine) startGame() {
	// Send each player info about its cards influences
	engine.sendCardInfluences()

	// Shuffle the players list
	rand.Seed(time.Now().UnixNano())
	players := engine.Game.Players
	rand.Shuffle(len(players), func(i, j int) { players[i], players[j] = players[j], players[i] })

	// Assign each player its gamePosition
	// and set the first player to act
	for i := 0; i < len(players); i++ {
		players[i].GamePosition = i
	}
	engine.Game.CurrentPlayer = players[0]

	engine.GlobalBroadcast(GameStarted)
}

func (engine *GameEngine) takeCoins(playerName string, coinsAmount int) {
	if engine.Game.TableCoins < coinsAmount {
		log.Fatal("Not enough coins on the table")
	}

	for i := 0; i < len(engine.Game.Players); i++ {
		player := &engine.Game.Players[i]
		if player.Name == playerName {
			player.Coins += coinsAmount
			engine.Game.TableCoins -= coinsAmount
			return
		}
	}
}

func (engine *GameEngine) waitForCounters() {
	engine.waitingCountersTimer = time.AfterFunc(4*time.Second, func() {
		engine.nextPlayer()
	})
}

func (engine *GameEngine) nextPlayer() {
	currentPosition := engine.Game.CurrentPlayer.GamePosition
	numPlayers := len(engine.Game.Players)

	currentPosition = (currentPosition + 1) % numPlayers
	engine.Game.CurrentPlayer = engine.Game.Players[currentPosition]

	engine.GlobalBroadcast(NextPlayer)
}

func (engine *GameEngine) solveChallenge(challengeType MessageType) {
	var challenged *models.Player
	var pretendingInfluence models.Influence
	var challengeSuccess bool

	game := engine.Game

	if challengeType == ChallengeAction {
		challenged = &game.CurrentPlayer

		switch game.CurrentPlayerAction.Action.ActionType {
		case models.TakeThreeCoins:
			pretendingInfluence = models.Duke
		case models.Exchange:
			pretendingInfluence = models.Ambassador
		case models.Steal:
			pretendingInfluence = models.Captain
		case models.Assasinate:
			pretendingInfluence = models.Assassin
		}
	} else if challengeType == ChallengeBlock {
		challenged = game.CurrentPlayerAction.BlockAction.Player
		pretendingInfluence = *game.CurrentPlayerAction.BlockAction.PretendingInfluence
	}

	// Check if the challenged player really has the pretending card
	if challenged.Card1.GetInfluence() == pretendingInfluence {
		challengeSuccess = false
		challenged.Card1 = game.InsertAndDraw(challenged.Card1)
	}
	if challenged.Card2.GetInfluence() == pretendingInfluence {
		challengeSuccess = false
		challenged.Card2 = game.InsertAndDraw(challenged.Card2)
	}

	// Send the challenged players its new card
	if !challengeSuccess {
		engine.sendPlayerCardInfluences(challenged)
	}

	// Update and broadcast the challenge result
	if challengeType == ChallengeAction {
		game.CurrentPlayerAction.ChallengeSuccess = &challengeSuccess
		engine.GlobalBroadcast(ChallenegeActionResult)
	}
	if challengeType == ChallengeBlock {
		game.CurrentPlayerAction.BlockAction.ChallengeSuccess = &challengeSuccess
		engine.GlobalBroadcast(ChallengeBlockResult)
	}
}
