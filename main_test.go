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

func sendCurrentPlayerAction(messageType core.MessageType, playerAction models.PlayerAction, engine *core.GameEngine) {
	sendPlayerAction(*engine.Game.CurrentPlayer, messageType, playerAction, engine)
}

func sendPlayerAction(player models.Player, messageType core.MessageType, playerAction models.PlayerAction, engine *core.GameEngine) {
	data, _ := json.Marshal(playerAction)

	payload, _ := json.Marshal(core.GameMessage{
		MessageType: messageType,
		Data:        data,
	})

	engine.ReadClientMessage(core.ClientMessage{
		ClientUuid: player.GetConnectionUuid(),
		Payload:    &payload,
	})

	time.Sleep(1 * time.Second)
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

	t.Run("TakeOneCoin action", func(t *testing.T) {
		playerAction := models.PlayerAction{
			Action: *models.NewAction(models.TakeOneCoin),
		}
		currentPlayer := engine.Game.CurrentPlayer
		sendCurrentPlayerAction(core.Action, playerAction, engine)

		if engine.Game.TableCoins != models.TOTAL_COINS-core.INITIAL_COINS_COUNT*core.MAX_PLAYERS-1 {
			t.Errorf("Wrong table coins after taking one coin. Expected : %d; Actual : %d",
				models.TOTAL_COINS-core.INITIAL_COINS_COUNT*core.MAX_PLAYERS-1,
				engine.Game.TableCoins)
		}
		if currentPlayer.Coins != 3 {
			t.Errorf("Wrong player coins after taking one coin. Expected : 3; Actual : %d", currentPlayer.Coins)
		}
	})

	t.Run("TakeTwoCoins action", func(t *testing.T) {
		currentPlayer := engine.Game.CurrentPlayer
		playerAction := models.PlayerAction{
			Action: *models.NewAction(models.TakeTwoCoins),
		}
		sendCurrentPlayerAction(core.Action, playerAction, engine)

		if currentPlayer.Coins != 4 {
			t.Errorf("Wrong player coins after taking one coin. Expected : 4; Actual : %d", currentPlayer.Coins)
		}
	})

	t.Run("Block with Duke", func(t *testing.T) {
		nextPlayer := engine.Game.Players[(engine.Game.CurrentPlayer.GamePosition+1)%core.MAX_PLAYERS]

		playerAction := engine.Game.CurrentPlayerAction
		influence := models.Duke
		playerAction.BlockAction = &models.Block{
			Player:              &nextPlayer,
			PretendingInfluence: &influence,
		}

		sendPlayerAction(nextPlayer, core.Block, *playerAction, engine)
	})

	t.Run("Challenge Block with Duke", func(t *testing.T) {
		playerAction := engine.Game.CurrentPlayerAction
		playerAction.BlockAction.ChallengedBy = engine.Game.CurrentPlayer

		sendCurrentPlayerAction(core.ChallengeBlock, *playerAction, engine)

		if playerAction.BlockAction.ChallengeSuccess == nil {
			t.Errorf("Wrong block with duke challenge result. Expected : true or false; Actual : nil")
		}

		if engine.Game.CurrentPlayer.HasInfluence(models.Duke) {
			if *playerAction.BlockAction.ChallengeSuccess == true {
				t.Error("Wrong block with duke challenge result. Expected : false; Actual : true")
			}
		} else {
			if *playerAction.BlockAction.ChallengeSuccess == false {
				t.Error("Wrong block with duke challenge result. Expected : true; Actual : false")
			}
		}
	})
}
