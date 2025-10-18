package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nathaponb/pon2go/api"
	"github.com/nathaponb/pon2go/internal/game"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	gameServer := game.NewServer()
	// start the game server main loop in goroutine
	go gameServer.Run(ctx)

	apiServer := api.NewApi(gameServer)
	go func() {
		// run api server in seperate goroutine
		log.Println("API server starting on port 3000")
		if err := apiServer.Listen(":3000"); err != nil {
			log.Printf("API server encountered error: %v", err)
		}
	}()

	go func() {
		// a goroutine listen on signal to shutdown api server
		<-ctx.Done() // Block until main's cancel() is called

		log.Println("API: Context canceled. Initiating Fiber Shutdown...")
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()

		if err := apiServer.ShutdownWithContext(shutdownCtx); err != nil {
			log.Printf("API: Fiber Shutdown Error: %v", err)
		} else {
			log.Println("API: Fiber shutdown complete.")
		}
	}()

	stop := make(chan os.Signal, 1) // size 1
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop // synchronization (block)

	log.Println("Received OS signal. Firing context cancellation...")

	cancel()

	time.Sleep(2 * time.Second)
	log.Println("Server exiting gracefully.")
}
