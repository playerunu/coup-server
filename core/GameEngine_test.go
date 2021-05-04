package core

import (
	"fmt"
	"testing"
)

func newTestEngine() *GameEngine {
	broadCastChannel := make(chan []byte)
	clientsPrivateChannel := make(chan ClientMessage)

	return NewGameEngine(&broadCastChannel, &clientsPrivateChannel)
}

func TestNewGame(t *testing.T) {
	engine := newTestEngine()

	fmt.Println("A trecut fratioare")
}
