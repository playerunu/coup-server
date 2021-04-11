package ws

import (
	"coup-server/core"
)

// GameServer maintains the set of active clients and broadcasts messages to the
// clients.
type GameServer struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcastChannel chan []byte

	// Register requests from the clients.
	registerChannel chan *Client

	// Unregister requests from clients.
	unregisterChannel chan *Client

	gameEngine *core.GameEngine
}

func NewGameServer() *GameServer {
	var gameServer = GameServer{
		broadcastChannel:  make(chan []byte),
		registerChannel:   make(chan *Client),
		unregisterChannel: make(chan *Client),
		clients:           make(map[*Client]bool),
	}

	gameServer.gameEngine = core.NewGameEngine(&gameServer.broadcastChannel)

	return &gameServer
}

func (gameServer *GameServer) Run() {
	go gameServer.gameEngine.Run()

	for {
		select {
		case client := <-gameServer.registerChannel:
			gameServer.clients[client] = true
		case client := <-gameServer.unregisterChannel:
			if _, ok := gameServer.clients[client]; ok {
				delete(gameServer.clients, client)
				close(client.sendChannel)
			}
		case message := <-gameServer.broadcastChannel:
			for client := range gameServer.clients {
				select {
				case client.sendChannel <- message:
				default:
					close(client.sendChannel)
					delete(gameServer.clients, client)
				}
			}
		}
	}
}
