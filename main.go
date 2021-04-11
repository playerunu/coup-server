package main

import (
	"coup-server/ws"
	"fmt"
	"log"
	"net/http"
)

const WS_SERVER_PORT = 10002

func initWsRoutes() {
	gameServer := ws.NewGameServer()
	go gameServer.Run()

	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		ws.OnWsConnect(gameServer, writer, request)
	})

	err := http.ListenAndServe(":"+fmt.Sprintf("%d", WS_SERVER_PORT), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func main() {
	fmt.Println("Starting server")
	initWsRoutes()
}
