package main

import (
	"github.com/nathaponb/pon2go/api"
	"github.com/nathaponb/pon2go/internal/game"
)

func main() {
	gameServer := game.NewServer()
	// start the game server main loop in goroutine
	go gameServer.Run()

	err := api.NewApi(gameServer)
	if err != nil {
		panic(err)
	}
}
