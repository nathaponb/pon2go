package game

import "github.com/gofiber/contrib/websocket"

type Player struct {
	Server *Server
	Conn   *websocket.Conn
	Send   chan []byte
	ID     string
}
