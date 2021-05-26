package ws

import (
	"coup-server/core"
)

// GameServer maintains the set of active clients and broadcasts messages to the
// clients.
type GameServer struct {
	// Registered clients.
	clients map[*GameClient]bool

	// Broadcast to all clients
	broadcastChannel chan []byte

	// Send private message to clients
	clientsPrivateChannel chan core.ClientMessage

	// Register connect requests from the clients.
	registerChannel chan *GameClient

	// Notifies the server about the player join / game startedupdates
	registerConfirmationChannel chan bool

	// Unregister disconnect requests from clients.
	unregisterChannel chan *GameClient

	gameEngine *core.GameEngine
}

func NewGameServer() *GameServer {
	var gameServer = GameServer{
		broadcastChannel:            make(chan []byte),
		registerChannel:             make(chan *GameClient),
		registerConfirmationChannel: make(chan bool),
		unregisterChannel:           make(chan *GameClient),
		clientsPrivateChannel:       make(chan core.ClientMessage),
		clients:                     make(map[*GameClient]bool),
	}

	gameServer.gameEngine = core.NewGameEngine(&gameServer.broadcastChannel, &gameServer.clientsPrivateChannel)

	return &gameServer
}

func (gameServer *GameServer) Run() {
	go gameServer.gameEngine.Run()

	for {
		select {
		case client := <-gameServer.registerChannel:
			gameServer.clients[client] = true
			if len(gameServer.clients) == core.MAX_PLAYERS {
				// Game sever is not full yet
				gameServer.registerConfirmationChannel <- true
			} else {
				// Game server is full
				gameServer.registerConfirmationChannel <- false
			}
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
		case message := <-gameServer.clientsPrivateChannel:
			for client := range gameServer.clients {
				if client.connectionUuid == message.ClientUuid {
					client.sendChannel <- *message.Payload
					break
				}
			}
		}
	}
}
