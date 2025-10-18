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
			s.matchPlayers(player)
		}
	}
}

func (s *Server) matchPlayers(newPlayer *Player) {
	// Check if a player is in waiting pool
	if len(s.waitingPool) > 0 {
		// Get the first player
		var waitingPlayer *Player
		for player := range s.waitingPool {
			waitingPlayer = player
			break
		}

		// Remove the waiting player from the waiting pool
		delete(s.waitingPool, waitingPlayer)

		log.Printf("Match found! Players: %s vs %s", waitingPlayer.ID, newPlayer.ID)

		// Create new room
		roomID := generateUniqueID()
		newRoom := NewRoom(roomID, waitingPlayer, newPlayer)

		// Map new room to game server
		s.room[roomID] = newRoom

		// Assign room back to players
		waitingPlayer.Room = newRoom
		newPlayer.Room = newRoom

		// Start the game in another goroutine
		go newRoom.Run()
	} else {
		// No active player in the waiting pool
		s.waitingPool[newPlayer] = true
		log.Printf("Player %s added to waiting pool. Current size: %d", newPlayer.ID, len(s.waitingPool))

		// Notify the player they are waiting
		newPlayer.Send <- []byte(`{"type": "WAITING", "message": "Waiting for an opponent..."}`)
	}
}

func (s *Server) cleanRooms() {}

func generateUniqueID() string {
	return ""
}
