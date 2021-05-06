package main

import (
	"coup-server/core"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func newTestEngine() *core.GameEngine {
	broadCastChannel := make(chan []byte)
	clientsPrivateChannel := make(chan core.ClientMessage)

	return core.NewGameEngine(&broadCastChannel, &clientsPrivateChannel)
}

func registerPlayer(playerName string, engine *core.GameEngine) {
	uuid := uuid.New()

	data, err := json.Marshal(struct {
		Name string
	}{
		Name: playerName,
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
		ClientUuid: uuid,
		Payload:    &payload,
	})
}

func TestNewGame(t *testing.T) {
	var playerNames = []string{"SerifIntergalactic", "NuSuntBot", "Capetanos", "DucuBertzi"}

	engine := newTestEngine()
	go engine.Run()

	t.Run("Initial game state", func(t *testing.T) {
		if engine.Game.TableCoins != 50 {
			t.Errorf("Wrong intial table coins count. Expected : 50; Actual : %d", engine.Game.TableCoins)
		}
	})

	t.Run("Register players", func(t *testing.T) {
		for index := 0; index < core.MAX_PLAYERS; index++ {
			registerPlayer(playerNames[index], engine)
		}

		if len(engine.Game.Players) != core.MAX_PLAYERS {
			t.Errorf("Wrong players count. Expected : %d; Actual : %d", core.MAX_PLAYERS, len(engine.Game.Players))
		}

		for index, player := range engine.Game.Players {
			if player.Name != playerNames[index] {
				t.Errorf("Wrong player name. Expected : %s; Actual : %s", playerNames[index], player.Name)
			}
			if player.Coins != core.INITIAL_COINS_COUNT {
				t.Errorf("Wrong intial coins count. Expected : %d; Actual : %d", core.INITIAL_COINS_COUNT, player.Coins)
			}
		}
	})
}
