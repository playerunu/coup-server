package main

import (
	"coup-server/core"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
)

func newTestEngine() *core.GameEngine {
	broadCastChannel := make(chan []byte)
	clientsPrivateChannel := make(chan core.ClientMessage)

	return core.NewGameEngine(&broadCastChannel, &clientsPrivateChannel)
}

func TestNewGame(t *testing.T) {
	engine := newTestEngine()
	if engine.Game.TableCoins != 50 {
		t.Errorf("Wrong intial table coins count. Expected : 50; Actual : %d", engine.Game.TableCoins)
	}

	go engine.Run()

	// Create 4 fake clients and register them to the game
	uuid1 := uuid.New()
	// uuid2 := uuid.New()
	// uuid3 := uuid.New()
	// uuid4 := uuid.New()

	data, err := json.Marshal(struct {
		Name string
	}{
		Name: "Serif",
	})

	if err != nil {
		fmt.Println("ceva")
	}
	payload, err := json.Marshal(core.GameMessage{
		MessageType: core.PlayerJoined,
		Data:        data,
	})

	if err != nil {
		fmt.Println("ceva")
	}

	engine.ReadClientMessage(core.ClientMessage{
		ClientUuid: uuid1,
		Payload:    &payload,
	})

	time.Sleep(1 * time.Second)

	if len(engine.Game.Players) != 1 {
		t.Errorf("Wrong players count. Expected : 1; Actual : %d", len(engine.Game.Players))
	}

	fmt.Println("A trecut fratioare")
}
