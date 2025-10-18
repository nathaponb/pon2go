package game

import "time"

type Room struct {
	id                 string
	player1            *Player
	player2            *Player
	ballX, ballY       float64
	paddle1Y, paddle2Y float64
	score1, score2     int
	Input              chan PlayerInput
	ticker             *time.Ticker
}

func NewRoom(id string, player1, player2 *Player) *Room {
	return &Room{
		id:      id,
		player1: player1,
		player2: player2,
	}
}

func (r *Room) Run() {}

func (r *Room) broardcastState() {}
