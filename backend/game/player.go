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

func CreatePlayer(role int, conn *websocket.Conn) Player {
	return Player{
		Name:  "yawdat",
		Role:  role,
		Score: 999,
		Conn:  conn,
	}
}
