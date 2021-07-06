package main

import (
	"coup-server/ws"
	"fmt"
)

func main() {
	// test
	fmt.Println("Starting server")
	wsServer := ws.NewWsServer()
	wsServer.Run()
}
