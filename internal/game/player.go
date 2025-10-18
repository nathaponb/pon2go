package game

import (
	"encoding/json"
	"log"

	"github.com/gofiber/contrib/websocket"
)

type Player struct {
	Server *Server
	Conn   *websocket.Conn
	Send   chan []byte
	ID     string
	Room   *Room
}

func (p *Player) ReadPump() {
	defer func() {
		p.Server.deRegister <- p
		p.Conn.Close()
	}()

	for {
		_, message, err := p.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}
			break
		}
		log.Printf("Player %s received: %s", p.ID, message)

		// parse player input into room
		var inputCommand struct {
			Action string `json:"action"` // e.g., "move"
			Value  string `json:"value"`  // e.g., "up"
		}

		if err := json.Unmarshal(message, &inputCommand); err != nil {
			log.Printf("Player %s: failed to unmarshal input: %v", p.ID, err)
			continue // Skip this message and wait for the next one
		}

		playerInput := PlayerInput{
			PlayerID: p.ID,
			Command:  inputCommand.Value, // Use the value as the command for simplicity
		}

		// Check if the player is in an active game room
		if p.Room != nil {
			// Send the input command to the room's dedicated input channel
			// This is non-blocking because the Room.Input channel should be buffered.
			select {
			case p.Room.Input <- playerInput:
				// Successfully sent input to the room
			default:
				// The room's input channel is full (game loop is too slow)
				log.Printf("Room %s: input channel full, dropping command from %s", p.Room.id, p.ID)
			}
		} else {
			// If not in a room, handle server-level commands (like "match me")
			if inputCommand.Action == "match" {
				p.Server.requestMatch <- p
			}
		}
	}
}

func (p *Player) WritePump() {
	defer func() {
		p.Conn.Close()
	}()

	for message := range p.Send {
		if err := p.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("error writing message to %s: %v", p.ID, err)
			return
		}
	}
}
