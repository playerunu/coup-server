package main

import (
	"coup-server/ws"
	"fmt"
	"log"
	"net/http"
)

const WS_SERVER_PORT = 8080

func initWsRoutes() {
	hub := ws.NewHub()
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
	initWsRoutes()
}
