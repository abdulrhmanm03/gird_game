package game

import (
	"github.com/gorilla/websocket"
)

type Player struct {
	Name  string
	Role  int
	Score int
	Conn  *websocket.Conn
}

const initScore = 100

// NOTE: player name is hard coded for now
func CreatePlayer(role int, conn *websocket.Conn) Player {
	return Player{
		Name:  "yawdat",
		Role:  role,
		Score: initScore,
		Conn:  conn,
	}
}
