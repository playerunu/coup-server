package main

import (
	"coup-server/core"
	models "coup-server/model"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
)

func newTestEngine() *core.GameEngine {
	broadCastChannel := make(chan []byte)
	clientsPrivateChannel := make(chan core.ClientMessage)

	go func() {
		for {
			select {
			case message := <-broadCastChannel:
				_ = message
			case message := <-clientsPrivateChannel:
				_ = message
			}
		}
	}()

	return core.NewGameEngine(&broadCastChannel, &clientsPrivateChannel)
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
}

func sendPlayerAction(player models.Player, messageType core.MessageType, playerAction models.PlayerMove, engine *core.GameEngine) {
	data, _ := json.Marshal(playerAction)

	payload, _ := json.Marshal(core.GameMessage{
		MessageType: messageType,
		Data:        data,
	})

	engine.ReadClientMessage(core.ClientMessage{
		ClientUuid: player.GetConnectionUuid(),
		Payload:    &payload,
	})

	time.Sleep(200 * time.Millisecond)
}

func sendRandomCardReveal(player models.Player, engine *core.GameEngine) {
	if !player.Card1.IsRevealed {
		player.Card1.Reveal()
	} else if !player.Card2.IsRevealed {
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

	time.Sleep(200 * time.Millisecond)
}

func getNextPlayer(engine *core.GameEngine) *models.Player {
	return &engine.Game.Players[(engine.Game.CurrentPlayer.GamePosition+1)%core.MAX_PLAYERS]
}

func TestGame(t *testing.T) {
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

	engine := newTestEngine()
	go engine.Run()

	t.Run("Initial game state", func(t *testing.T) {
		if engine.Game.TableCoins != models.TOTAL_COINS {
			t.Errorf("Wrong intial table coins count. Expected : %d; Actual : %d", models.TOTAL_COINS, engine.Game.TableCoins)
		}
	})

	t.Run("Register players", func(t *testing.T) {
		for index := 0; index < core.MAX_PLAYERS; index++ {
			registerPlayer(players[index], engine)
		}

		time.Sleep(1 * time.Second)

		if len(engine.Game.Players) != core.MAX_PLAYERS {
			t.Errorf("Wrong players count. Expected : %d; Actual : %d", core.MAX_PLAYERS, len(engine.Game.Players))
		}

		for _, player := range engine.Game.Players {
			if player.Coins != core.INITIAL_COINS_COUNT {
				t.Errorf("Wrong intial coins count. Expected : %d; Actual : %d", core.INITIAL_COINS_COUNT, player.Coins)
			}
		}
	})

	firstPlayer := engine.Game.CurrentPlayer
	secondPlayer := getNextPlayer(engine)

	t.Run("[Round 1] First player - TakeOneCoin", func(t *testing.T) {
		playerAction := models.PlayerMove{
			Action: *models.NewAction(models.TakeOneCoin),
		}
		sendPlayerAction(*firstPlayer, core.Action, playerAction, engine)

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
		sendPlayerAction(*secondPlayer, core.Action, playerMove, engine)

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

		sendPlayerAction(*firstPlayer, core.Block, *playerMove, engine)
	})

	t.Run("[Round 2] Second player - Challenges Block with Duke", func(t *testing.T) {
		playerMove := engine.Game.CurrentMove
		playerMove.Block.Challenge = &models.Challenge{
			ChallengedBy: secondPlayer,
		}

		hasDuke := firstPlayer.HasInfluence(models.Duke)
		sendPlayerAction(*secondPlayer, core.ChallengeBlock, *playerMove, engine)

		if playerMove.Block.Challenge == nil {
			t.Errorf("Wrong block with duke challenge result. Expected : a Block struct instance; Actual : nil")
		}

		if hasDuke {
			if *playerMove.Block.Challenge.Success == true {
				t.Error("Wrong block with duke challenge result. Expected : false; Actual : true")
			}

			sendRandomCardReveal(*firstPlayer, engine)
		} else {
			if *playerMove.Block.Challenge.Success == false {
				t.Error("Wrong block with duke challenge result. Expected : true; Actual : false")
			}

			sendRandomCardReveal(*secondPlayer, engine)
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
		sendPlayerAction(*secondPlayer, core.Action, playerMove, engine)

		// Wait and dont challenge the action
		time.Sleep(4 * time.Second)

		remainingCards := secondPlayer.RemainingCards()
		sendRandomCardReveal(*secondPlayer, engine)

		if remainingCards == 1 {
			if !secondPlayer.IsEliminated() {
				t.Errorf("Wrong player remaining cards. Expected: 0; Actual: %d", secondPlayer.RemainingCards())
			}
		} else {
			if secondPlayer.RemainingCards() != 1 {
				t.Errorf("Wrong player remaining cards. Expected: 1; Actual: %d", secondPlayer.RemainingCards())
			}
		}
	})
}
