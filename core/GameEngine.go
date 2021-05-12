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
const WAITING_COUNTERS_SECONDS = 3

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
			engine.onPlayerMove(gameMessage, clientMessage.ClientUuid)
		case ChallengeAction:
		case Block:
		case ChallengeBlock:
			engine.onPlayerCounter(gameMessage, clientMessage.ClientUuid)
		case RevealCard:
			engine.onCardReveal(gameMessage, clientMessage.ClientUuid)
		}
	}
}

func (engine *GameEngine) ReadClientMessage(message ClientMessage) {
	engine.clientUpdatesChannel <- message
}

func (engine *GameEngine) Broadcast(messageType MessageType) {
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
		log.Fatalln("Error while unmarshalling player join message:", err)
	}
	player.SetConnectionUuid(uuid)

	engine.registerPlayer(player)
}

func (engine *GameEngine) onPlayerMove(message GameMessage, uuid uuid.UUID) {
	var clientAction models.PlayerMove
	err := json.Unmarshal(message.Data, &clientAction)
	if err != nil {
		log.Fatalln("Error while unmarshalling game message: ", message.MessageType, err)
	}

	game := engine.Game
	var playerAction models.PlayerMove
	playerAction.Action = *models.NewAction(clientAction.Action.ActionType)

	switch clientAction.Action.ActionType {
	case models.TakeOneCoin:
		engine.getCoinsFromTable(game.CurrentPlayer, 1)
	case models.TakeTwoCoins:
		engine.getCoinsFromTable(game.CurrentPlayer, 2)
	case models.TakeThreeCoins:
		engine.getCoinsFromTable(game.CurrentPlayer, 3)
	case models.Assasinate:
	case models.Steal:
		playerAction.VsPlayer = clientAction.VsPlayer
	case models.Exchange:
		// This action will only be broadcasted
		break
	}

	engine.Game.CurrentMove = &playerAction
	engine.Broadcast(Action)

	if playerAction.CanCounter() {
		engine.waitForCounters()
	} else {
		engine.finishCurrentMove()
	}
}

func (engine *GameEngine) onPlayerCounter(message GameMessage, uuid uuid.UUID) {
	// Always make sure we first stop the waiting counters timer
	engine.waitingCountersTimer.Stop()

	var clientAction models.PlayerMove
	err := json.Unmarshal(message.Data, &clientAction)
	if err != nil {
		log.Fatalln("Error while unmarshalling game message: ", message.MessageType, err)
	}

	currentMove := engine.Game.CurrentMove

	switch message.MessageType {
	case ChallengeAction:
		currentMove.Challenge = &models.Challenge{
			ChallengedBy:  clientAction.Challenge.ChallengedBy,
			WaitingReveal: true,
		}
		engine.solveChallenge(message.MessageType)
	case ChallengeBlock:
		currentMove.Block.Challenge = &models.Challenge{
			ChallengedBy:  clientAction.Block.Challenge.ChallengedBy,
			WaitingReveal: true,
		}
		engine.solveChallenge(message.MessageType)
	case Block:
		currentMove.Block = &models.Block{
			Player:              clientAction.Block.Player,
			PretendingInfluence: clientAction.Block.PretendingInfluence,
		}
	}

	engine.Broadcast(message.MessageType)

	// Restart the waiting counters timer only when blocking
	// For challenges, we restart the counter after the card is revealed
	if message.MessageType == Block {
		engine.waitForCounters()
	}
}

func (engine *GameEngine) onCardReveal(message GameMessage, uuid uuid.UUID) {
	var playerUpdate models.Player
	err := json.Unmarshal(message.Data, &playerUpdate)
	if err != nil {
		log.Fatalln("Error while unmarshalling card reveal message:", err)
	}

	player := engine.Game.GetPlayerByName(playerUpdate.Name)
	if playerUpdate.Card1.IsRevealed {
		player.Card1.Reveal()
	}
	if playerUpdate.Card2.IsRevealed {
		player.Card2.Reveal()
	}

	currentMove := engine.Game.CurrentMove
	if currentMove.WaitingReveal != nil {
		*currentMove.WaitingReveal = false
	} else if currentMove.Challenge != nil && currentMove.Challenge.WaitingReveal {
		currentMove.Challenge.WaitingReveal = false
	} else if currentMove.Block.Challenge != nil && currentMove.Block.Challenge.WaitingReveal {
		currentMove.Block.Challenge.WaitingReveal = false
	}

	engine.Broadcast(message.MessageType)

	if engine.Game.CurrentMove.CanCounter() {
		engine.waitForCounters()
	} else {
		engine.finishCurrentMove()
	}
}

// Registers a new player
func (engine *GameEngine) registerPlayer(player models.Player) {
	// Draw 2 random cards
	player.Card1 = engine.Game.DrawCard()
	player.Card2 = engine.Game.DrawCard()

	// Give the initial coins
	engine.getCoinsFromTable(&player, INITIAL_COINS_COUNT)

	engine.Game.Players = append(engine.Game.Players, player)

	if len(engine.Game.Players) == MAX_PLAYERS {
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
	engine.Game.CurrentPlayer = &players[0]

	engine.Broadcast(GameStarted)
}

func (engine *GameEngine) getCoinsFromTable(player *models.Player, coinsAmount int) {
	if engine.Game.TableCoins < coinsAmount {
		log.Fatal("Not enough coins on the table")
	}

	player.Coins += coinsAmount
	engine.Game.TableCoins -= coinsAmount
}

func (engine *GameEngine) putCoinsOnTable(player *models.Player, coinsAmount int) {
	player.Coins -= coinsAmount
	engine.Game.TableCoins += coinsAmount
}

func (engine *GameEngine) waitForCounters() {
	engine.waitingCountersTimer = time.AfterFunc(WAITING_COUNTERS_SECONDS*time.Second, func() {
		engine.finishCurrentMove()
	})
}

func (engine *GameEngine) finishCurrentMove() {
	currentMove := engine.Game.CurrentMove
	currentPlayer := engine.Game.CurrentPlayer
	vsPlayer := engine.Game.CurrentMove.VsPlayer
	waitingReveal := false
	waitingExchange := false

	// Finalize the current action result
	if currentMove.IsSuccessful() {
		switch currentMove.Action.ActionType {
		case models.Assasinate:
			if currentMove.WaitingReveal == nil || !*currentMove.WaitingReveal {
				engine.putCoinsOnTable(currentPlayer, models.ASSASSINATE_COINS_AMOUNT)
				waitingReveal = true
				currentMove.WaitingReveal = &waitingReveal
			}
		case models.Coup:
			engine.putCoinsOnTable(currentPlayer, models.COUP_COINS_AMOUNT)
		case models.Steal:
			if vsPlayer.Coins >= 2 {
				currentPlayer.Coins += 2
				vsPlayer.Coins -= 2
			} else {
				currentPlayer.Coins = vsPlayer.Coins
				vsPlayer.Coins = 0
			}
		case models.Exchange:
			if currentMove.WaitingExchange == nil || !*currentMove.WaitingExchange {
				waitingExchange = true
				currentMove.WaitingExchange = &waitingExchange
			}
		}
	} else {
		// Some actions need to rollback when they are blocked/challenged
		switch currentMove.Action.ActionType {
		case models.TakeTwoCoins:
			engine.putCoinsOnTable(currentPlayer, 2)
		case models.TakeThreeCoins:
			engine.putCoinsOnTable(currentPlayer, 3)
		case models.Assasinate:
			engine.getCoinsFromTable(currentPlayer, 3)
		}
	}

	engine.Broadcast(ActionResult)
	// We still have to wait for the assassinate/coup/exchange to happen before moving on
	if !waitingExchange && !waitingReveal {
		engine.nextPlayer()
	}
}

func (engine *GameEngine) nextPlayer() {
	// Calculate the next player position
	currentPosition := engine.Game.CurrentPlayer.GamePosition
	numPlayers := len(engine.Game.Players)

	currentPosition = (currentPosition + 1) % numPlayers
	engine.Game.CurrentPlayer = &engine.Game.Players[currentPosition]
	engine.Game.CurrentMove = nil

	engine.Broadcast(NextPlayer)
}

func (engine *GameEngine) solveChallenge(challengeType MessageType) {
	var challenged *models.Player
	var challenger *models.Player
	var challenge *models.Challenge
	var pretendingInfluence *models.Influence
	var success bool

	game := engine.Game
	currentMove := game.CurrentMove

	if challengeType == ChallengeAction {
		challenged = game.CurrentPlayer
		challenger = currentMove.Challenge.ChallengedBy
		challenge = currentMove.Challenge
		pretendingInfluence = currentMove.Action.GetInfluence()
	} else if challengeType == ChallengeBlock {
		challenged = currentMove.Block.Player
		challenger = currentMove.Block.Challenge.ChallengedBy
		challenge = currentMove.Block.Challenge
		pretendingInfluence = currentMove.Block.PretendingInfluence
	}

	// Check if the challenged player really has the pretending card
	if challenged.Card1.GetInfluence() == *pretendingInfluence {
		success = false
		challenged.Card1 = game.InsertAndDraw(challenged.Card1)
	} else if challenged.Card2.GetInfluence() == *pretendingInfluence {
		success = false
		challenged.Card2 = game.InsertAndDraw(challenged.Card2)
	} else {
		success = true
	}

	if success {
		// Challenge won, challenged player should reveal a card
		if challenged.RemainingCards() == 1 {
			challenged.RevealLastCard()
		}
	} else {
		// Challenge lost, challenger player should reveal a card
		if challenger.RemainingCards() == 1 {
			challenger.RevealLastCard()
		}

		// Send the challenged players its new card
		engine.sendPlayerCardInfluences(challenged)
	}

	// Broadcast the challenge result
	challenge.Success = &success
	if challengeType == ChallengeAction {
		engine.Broadcast(ChallenegeActionResult)
	} else if challengeType == ChallengeBlock {
		engine.Broadcast(ChallengeBlockResult)
	}

	// If the action can still be countered, wait for the card reveal first,
	// then start the waitCounter timer, otherwise
	// move to next player
	if !challenge.WaitingReveal && !currentMove.CanCounter() {
		engine.finishCurrentMove()
	}
}
