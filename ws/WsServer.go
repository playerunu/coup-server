package ws

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
)

const WS_SERVER_PORT = 10002

type WsServer struct {
	runningGames     []*GameServer
	currentGame      *GameServer
	currentGameMutex sync.Mutex
}

var (
	once             sync.Once
	wsServerInstance *WsServer
)

func NewWsServer() *WsServer {
	once.Do(func() {
		wsServerInstance = &WsServer{
			currentGame: newGame(),
		}
	})

	return wsServerInstance
}

func (wsServer *WsServer) Run() {
	NewWsServer()

	gameServer := NewGameServer()
	go gameServer.Run()
	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		wsServer.OnWsConnect(writer, request)
	})

	err := http.ListenAndServe(":"+fmt.Sprintf("%d", WS_SERVER_PORT), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

// serveWs handles websocket requests from the peer.
func (wsServer *WsServer) OnWsConnect(w http.ResponseWriter, r *http.Request) {

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	wsServer.currentGameMutex.Lock()

	gameServer := wsServer.currentGame
	client := &GameClient{
		connectionUuid:    uuid.New(),
		unregisterChannel: &gameServer.unregisterChannel,
		gameEngineChannel: &gameServer.gameEngine.ClientUpdatesChannel,
		conn:              conn,
		sendChannel:       make(chan []byte, 256),
	}

	go client.runReader()
	go client.runWriter()

	gameServer.registerChannel <- client

	confirmation := <-gameServer.registerConfirmationChannel
	if confirmation {
		wsServer.runningGames = append(wsServer.runningGames, gameServer)
		wsServer.currentGame = newGame()
	}

	wsServer.currentGameMutex.Unlock()
}

func newGame() *GameServer {
	gameServer := NewGameServer()
	go gameServer.Run()
	return gameServer
}
