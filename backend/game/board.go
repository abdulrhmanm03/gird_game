package game

import (
	"log"
	"time"
)

type Squere struct {
	Content int
	Clicked chan bool
}

func newSquere() *Squere {
	return &Squere{
		Content: 0,
		Clicked: make(chan bool),
	}
}

type Msg struct {
	Pos   int `json:"pos"`
	Score int `json:"score"`
}

func (s *Squere) Run(player *Player, index int) {
	select {
	case <-s.Clicked:
		break
	case <-time.After(15 * time.Second):
		if s.Content == 1 {
			player.Score -= 5
		} else if s.Content == 2 {
			player.Score += 5
		}
		msg := Msg{
			Pos:   index,
			Score: player.Score,
		}
		err := player.Conn.WriteJSON(msg)
		if err != nil {
			log.Println("write faild to write to player")
		}
	}
	s.Content = 0
}

func CreateBoard(width int, height int) []*Squere {
	board := make([]*Squere, width*height)
	for i := range board {
		board[i] = newSquere()
	}
	return board
}
