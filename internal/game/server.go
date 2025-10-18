package game

import (
	"context"
	"log"
)

type Server struct {
	waitingPool  map[*Player]bool
	room         map[string]*Room
	Register     chan *Player
	deRegister   chan *Player
	requestMatch chan *Player
}

func NewServer() *Server {
	return &Server{
		waitingPool:  make(map[*Player]bool),
		room:         make(map[string]*Room),
		Register:     nil,
		deRegister:   nil,
		requestMatch: nil,
	}
}

func (s *Server) Run(ctx context.Context) {
	log.Println("Game server is running...")
	for {
		select {
		case <-ctx.Done():
			log.Println("Receive context cancellation, cleaning up resources and shut down game server")
			// cleanup resources
			s.cleanRooms()
			return // terminates the Run() function and its goroutine.
		case player := <-s.Register:
			log.Printf("Platyer %s registered", player.ID)
			// send new register player to requestMatch channel
			s.requestMatch <- player
		case player := <-s.deRegister:
			log.Printf("Platyer %s de-registered", player.ID)
			// clean up player and potentially room
			delete(s.waitingPool, player)
			// close player send channel
			close(player.Send)
		case player := <-s.requestMatch:
			log.Printf("Platyer %s requested match, Pool size: %d", player.ID, len(s.waitingPool))

			//TODO: implement match-making logic

		}
	}
}

func (s *Server) cleanRooms() {}
