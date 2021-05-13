package main

import (
	"coup-server/core"
	models "coup-server/model"
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/google/uuid"
)

var players = []models.Player{
	{
		Name: "SerifIntergalactic",
	},
	{
		Name: "NuSuntBot",
	},
	{
		Name: "Capetanos",
	},
	{
		Name: "DucuBertzi",
	},
}

var engine *core.GameEngine
var broadcast chan []byte
var clientUpdates chan core.ClientMessage
var firstPlayer *models.Player
var secondPlayer *models.Player

func newTestEngine() (*core.GameEngine, chan []byte, chan core.ClientMessage) {
	broadCastChannel := make(chan []byte)
	clientsPrivateChannel := make(chan core.ClientMessage)

	return core.NewGameEngine(&broadCastChannel, &clientsPrivateChannel), broadCastChannel, clientsPrivateChannel
}

func registerPlayer(player models.Player, engine *core.GameEngine) {
	player.SetConnectionUuid(uuid.New())

	data, err := json.Marshal(struct {
		Name string
	}{
		Name: player.Name,
	})

	if err != nil {
		fmt.Println("Error while marshalling PlayerJoin data", err)
	}

	payload, err := json.Marshal(core.GameMessage{
		MessageType: core.PlayerJoined,
		Data:        data,
	})

	if err != nil {
		fmt.Println("Error while marshalling PlayerJoin game message", err)
	}

	engine.ReadClientMessage(core.ClientMessage{
		ClientUuid: player.GetConnectionUuid(),
		Payload:    &payload,
	})

	time.Sleep(100 * time.Millisecond)
	if len(engine.Game.Players) == core.MAX_PLAYERS {
		// Game start message
		<-broadcast
	}
}

func sendPlayerMove(player models.Player, messageType core.MessageType, playerMove models.PlayerMove) {
	data, _ := json.Marshal(playerMove)

	payload, _ := json.Marshal(core.GameMessage{
		MessageType: messageType,
		Data:        data,
	})

	engine.ReadClientMessage(core.ClientMessage{
		ClientUuid: player.GetConnectionUuid(),
		Payload:    &payload,
	})
}

func sendActionWaitReveal(player models.Player, messageType core.MessageType, playerMove models.PlayerMove) {
	sendPlayerMove(player, messageType, playerMove)
	// Action message
	<-broadcast
	// Action result message
	<-broadcast
}

func sendActionNoCounter(player models.Player, messageType core.MessageType, playerMove models.PlayerMove) {
	sendPlayerMove(player, messageType, playerMove)
	// Action message
	<-broadcast
	// Action result message
	<-broadcast
	// Next player OR Game over message
	<-broadcast
}

func sendAction(player models.Player, messageType core.MessageType, playerMove models.PlayerMove) {
	sendPlayerMove(player, messageType, playerMove)
	// Action message
	<-broadcast
}

func sendReveal(player models.Player, cardType core.CardType) {
	if cardType == core.AnyUnrevealed {
		if !player.Card1.IsRevealed {
			cardType = core.Card1
		} else if !player.Card2.IsRevealed {
			cardType = core.Card2
		} else {
			log.Fatal("Both player cards are revealed and we are supposed to find one unrevealed card")
		}
	}

	if cardType == core.Card1 {
		player.Card1.Reveal()
	} else {
		player.Card2.Reveal()
	}

	data, _ := json.Marshal(player)

	payload, _ := json.Marshal(core.GameMessage{
		MessageType: core.RevealCard,
		Data:        data,
	})

	engine.ReadClientMessage(core.ClientMessage{
		ClientUuid: player.GetConnectionUuid(),
		Payload:    &payload,
	})

	// Reveal card
	<-broadcast
	// Next player
	<-broadcast
}

func getNextPlayer() *models.Player {
	return &engine.Game.Players[(engine.Game.CurrentPlayer.GamePosition+1)%core.MAX_PLAYERS]
}

func initTest() {
	engine, broadcast, clientUpdates = newTestEngine()
	go func() {
		for {
			<-clientUpdates
		}
	}()

	go engine.Run()

	for index := 0; index < core.MAX_PLAYERS; index++ {
		registerPlayer(players[index], engine)
	}

	firstPlayer = engine.Game.CurrentPlayer
	secondPlayer = getNextPlayer()
}

func TestMultipleRounds(t *testing.T) {
	initTest()

	t.Run("Initial game state", func(t *testing.T) {
		if engine.Game.TableCoins != models.TOTAL_COINS-core.MAX_PLAYERS*core.INITIAL_COINS_COUNT {
			t.Errorf("Wrong intial table coins count. Expected : %d; Actual : %d", models.TOTAL_COINS, engine.Game.TableCoins)
		}
	})

	t.Run("Register players", func(t *testing.T) {
		if len(engine.Game.Players) != core.MAX_PLAYERS {
			t.Errorf("Wrong players count. Expected : %d; Actual : %d", core.MAX_PLAYERS, len(engine.Game.Players))
		}

		for _, player := range engine.Game.Players {
			if player.Coins != core.INITIAL_COINS_COUNT {
				t.Errorf("Wrong intial coins count. Expected : %d; Actual : %d", core.INITIAL_COINS_COUNT, player.Coins)
			}
		}
	})

	t.Run("[Round 1] First player - TakeOneCoin", func(t *testing.T) {
		playerMove := models.PlayerMove{
			Action: *models.NewAction(models.TakeOneCoin),
		}
		sendActionNoCounter(*firstPlayer, core.Action, playerMove)

		if engine.Game.TableCoins != models.TOTAL_COINS-core.INITIAL_COINS_COUNT*core.MAX_PLAYERS-1 {
			t.Errorf("Wrong table coins after taking one coin. Expected : %d; Actual : %d",
				models.TOTAL_COINS-core.INITIAL_COINS_COUNT*core.MAX_PLAYERS-1,
				engine.Game.TableCoins)
		}
		if firstPlayer.Coins != 3 {
			t.Errorf("Wrong player coins after taking one coin. Expected : 3; Actual : %d", firstPlayer.Coins)
		}
	})

	t.Run("[Round 2] Second player - TakeTwoCoins", func(t *testing.T) {
		if engine.Game.CurrentPlayer != secondPlayer {
			t.Error("Wrong current player. Expected: second player")
		}
		playerMove := models.PlayerMove{
			Action: *models.NewAction(models.TakeTwoCoins),
		}
		sendAction(*secondPlayer, core.Action, playerMove)

		if secondPlayer.Coins != 4 {
			t.Errorf("Wrong player coins after taking two coins. Expected : 4; Actual : %d", secondPlayer.Coins)
		}
	})

	t.Run("[Round 2] First player - Block with Duke", func(t *testing.T) {
		playerMove := engine.Game.CurrentMove
		influence := models.Duke
		playerMove.Block = &models.Block{
			Player:              firstPlayer,
			PretendingInfluence: &influence,
		}

		sendAction(*firstPlayer, core.Block, *playerMove)
	})

	t.Run("[Round 2] Second player - Challenges Block with Duke", func(t *testing.T) {
		playerMove := engine.Game.CurrentMove
		playerMove.Block.Challenge = &models.Challenge{
			ChallengedBy: secondPlayer,
		}

		hasDuke := firstPlayer.HasInfluence(models.Duke)
		sendActionWaitReveal(*secondPlayer, core.ChallengeBlock, *playerMove)

		if playerMove.Block.Challenge == nil {
			t.Errorf("Wrong block with duke challenge result. Expected : a Block struct instance; Actual : nil")
		}

		if hasDuke {
			if *playerMove.Block.Challenge.Success == true {
				t.Error("Wrong block with duke challenge result. Expected : false; Actual : true")
			}

			sendReveal(*firstPlayer, core.AnyUnrevealed)
		} else {
			if *playerMove.Block.Challenge.Success == false {
				t.Error("Wrong block with duke challenge result. Expected : true; Actual : false")
			}

			sendReveal(*secondPlayer, core.AnyUnrevealed)
		}
	})

	t.Run("[Round 3] First player - Assassinate", func(t *testing.T) {
		if engine.Game.CurrentPlayer != firstPlayer {
			t.Errorf("Wrong current player. Expected: firstPlayer; Actual: secondPlayer")
		}

		playerMove := models.PlayerMove{
			Action:   *models.NewAction(models.Assasinate),
			VsPlayer: secondPlayer,
		}
		sendActionWaitReveal(*secondPlayer, core.Action, playerMove)

		remainingCards := secondPlayer.RemainingCards()
		sendReveal(*secondPlayer, core.AnyUnrevealed)

		if remainingCards == 1 {
			if !secondPlayer.IsEliminated() {
				t.Errorf("Wrong player remaining cards. Expected: 0; Actual: %d", secondPlayer.RemainingCards())
			}
			if engine.Game.Winner != firstPlayer {
				t.Errorf("Wrong game state - it should have a winner")
			}
		} else {
			if secondPlayer.RemainingCards() != 1 {
				t.Errorf("Wrong player remaining cards. Expected: 1; Actual: %d", secondPlayer.RemainingCards())
			}
		}
	})
}

func TestCoup(t *testing.T) {
	initTest()

	for i := 0; i < 7; i++ {
		playerMove := models.PlayerMove{
			Action: *models.NewAction(models.TakeOneCoin),
		}
		sendActionNoCounter(*firstPlayer, core.Action, playerMove)
		sendActionNoCounter(*secondPlayer, core.Action, playerMove)
	}

	playerMove := models.PlayerMove{
		Action:   *models.NewAction(models.Coup),
		VsPlayer: secondPlayer,
	}

	sendActionWaitReveal(*firstPlayer, core.Action, playerMove)
	sendReveal(*secondPlayer, core.AnyUnrevealed)

	if secondPlayer.RemainingCards() != 1 {
		t.Errorf("Wrong second player remaining cards. Expected : 1; Actual : %d", secondPlayer.RemainingCards())
	}

	playerMove = models.PlayerMove{
		Action:   *models.NewAction(models.Coup),
		VsPlayer: firstPlayer,
	}

	sendActionWaitReveal(*secondPlayer, core.Action, playerMove)
	sendReveal(*firstPlayer, core.AnyUnrevealed)

	if firstPlayer.RemainingCards() != 1 {
		t.Errorf("Wrong second player remaining cards. Expected : 1; Actual : %d", firstPlayer.RemainingCards())
	}

	for i := 0; i < 7; i++ {
		playerMove := models.PlayerMove{
			Action: *models.NewAction(models.TakeOneCoin),
		}
		sendActionNoCounter(*firstPlayer, core.Action, playerMove)
		sendActionNoCounter(*secondPlayer, core.Action, playerMove)
	}

	playerMove = models.PlayerMove{
		Action:   *models.NewAction(models.Coup),
		VsPlayer: secondPlayer,
	}

	sendActionWaitReveal(*firstPlayer, core.Action, playerMove)
	sendReveal(*secondPlayer, core.AnyUnrevealed)

	if secondPlayer.RemainingCards() != 0 {
		t.Errorf("Wrong second player remaining cards. Expected : 0; Actual : %d", secondPlayer.RemainingCards())
	}

	if engine.Game.Winner != firstPlayer {
		t.Errorf("Wrong player winner. It should be the first player")
	}
}

func TestSteal(t *testing.T) {
	var playerMove models.PlayerMove
	initTest()

	core.DrawInfluence(engine, firstPlayer, models.Captain, core.Card1)
	core.DrawInfluence(engine, firstPlayer, models.Contessa, core.Card2)

	core.DrawInfluence(engine, secondPlayer, models.Ambassador, core.Card1)
	core.DrawInfluence(engine, secondPlayer, models.Assassin, core.Card2)

	playerMove = *models.NewPlayerMove(models.Steal, secondPlayer)
	sendActionNoCounter(*firstPlayer, core.Action, playerMove)

	if secondPlayer.Coins != 0 {
		t.Errorf("Wrong second player coins. Expected : 0; Actual : %d", secondPlayer.Coins)
	}
	if engine.Game.CurrentPlayer != secondPlayer {
		t.Errorf("Wrong current player. Expected : second player")
	}

	playerMove = *models.NewPlayerMove(models.TakeOneCoin, nil)
	sendActionNoCounter(*secondPlayer, core.Action, playerMove)

	playerMove = *models.NewPlayerMove(models.Steal, secondPlayer)
	sendActionNoCounter(*firstPlayer, core.Action, playerMove)

	if secondPlayer.Coins != 0 {
		t.Errorf("Wrong second player coins. Expected : 0; Actual : %d", secondPlayer.Coins)
	}
	if firstPlayer.Coins != 5 {
		t.Errorf("Wrong first player coins. Expected : 5; Actual : %d", firstPlayer.Coins)
	}

	playerMove = *models.NewPlayerMove(models.TakeOneCoin, nil)
	sendActionNoCounter(*secondPlayer, core.Action, playerMove)

	playerMove = *models.NewPlayerMove(models.Steal, secondPlayer)
	sendAction(*firstPlayer, core.Action, playerMove)

	playerMove = *engine.Game.CurrentMove
	influence := models.Ambassador
	playerMove.Block = &models.Block{
		Player:              secondPlayer,
		PretendingInfluence: &influence,
	}

	sendActionNoCounter(*secondPlayer, core.Block, playerMove)

	if secondPlayer.Coins != 1 {
		t.Errorf("Wrong second player coins. Expected : 1; Actual : %d", secondPlayer.Coins)
	}

	playerMove = *models.NewPlayerMove(models.TakeOneCoin, nil)
	sendActionNoCounter(*secondPlayer, core.Action, playerMove)

	playerMove = *models.NewPlayerMove(models.Steal, secondPlayer)
	sendAction(*firstPlayer, core.Action, playerMove)

	playerMove = *engine.Game.CurrentMove
	influence = models.Ambassador
	playerMove.Block = &models.Block{
		Player:              secondPlayer,
		PretendingInfluence: &influence,
	}
	hasAmbassador := secondPlayer.HasInfluence(models.Ambassador)

	sendAction(*secondPlayer, core.Block, playerMove)

	playerMove = *engine.Game.CurrentMove
	playerMove.Block.Challenge = &models.Challenge{
		ChallengedBy: firstPlayer,
	}
	sendActionWaitReveal(*firstPlayer, core.ChallengeBlock, playerMove)

	if hasAmbassador {
		sendReveal(*firstPlayer, core.AnyUnrevealed)

		if firstPlayer.RemainingCards() != 1 {
			t.Errorf("Wrong remaining cards for first player. Expected : 1; Actual : %d", firstPlayer.RemainingCards())
		}
		if secondPlayer.Coins != 2 {
			t.Errorf("Wrong coins amount for second player. Expected : 2; Actual : %d", secondPlayer.Coins)
		}
	} else {
		sendReveal(*secondPlayer, core.AnyUnrevealed)

		if secondPlayer.RemainingCards() != 1 {
			t.Errorf("Wrong remaining cards for second player. Expected : 1; Actual : %d", secondPlayer.RemainingCards())
		}
		if secondPlayer.Coins != 0 {
			t.Errorf("Wrong coins amount for second player. Expected : 0; Actual : %d", secondPlayer.Coins)
		}
	}
}
