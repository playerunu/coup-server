package main

import (
	"coup-server/ws"
	"fmt"
)

func main() {
	fmt.Println("Starting server")
	wsServer := ws.NewWsServer()
	wsServer.Run()
}
