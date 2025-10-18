package api

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/nathaponb/pon2go/internal/game"
)

func NewApi(gs *game.Server) *fiber.App {
	app := fiber.New()

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/:id", websocket.New(func(c *websocket.Conn) {

		playerID := c.Params("id")
		player := &game.Player{
			ID:     playerID,
			Conn:   c,
			Server: gs,
			Send:   make(chan []byte, 256),
		}

		gs.Register <- player

	}))

	// return app.Listen(":3000")
	return app
}
