package game

import "time"

type Room struct {
	id                 string
	player1            *Player
	player2            *Player
	ballX, ballY       float64
	paddle1Y, paddle2Y float64
	score1, score2     int
	//input              chan PlayerInput
	ticker *time.Ticker
}

func (r *Room) run() {}

func (r *Room) broardcastState() {}
