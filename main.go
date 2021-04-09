package main

import (
	"coup-server/ws"
	"coup-server/models"
	"fmt"
	"log"
	"net/http"
)

const WS_SERVER_PORT = 10002

func initWsRoutes(game models.Game) {
	hub := ws.NewHub(game)
	go hub.Run()

	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		ws.ServeWs(hub, writer, request)
	})

	err := http.ListenAndServe(":"+fmt.Sprintf("%d", WS_SERVER_PORT), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func main() {
	fmt.Println("Starting server")
	game := models.NewGame()
	initWsRoutes(game)
}
