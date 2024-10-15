package game

import (
	"time"
)

const (
	Empty = iota
	Bomb
	Apple
)

type Cell struct {
	Content  int
	Clicked  chan bool
	TimeOver chan bool
}

func newCell() *Cell {
	return &Cell{
		Content:  Empty,
		Clicked:  make(chan bool),
		TimeOver: make(chan bool),
	}
}

func (s *Cell) Run(player *Player, pos int) {
	select {
	case <-s.Clicked:
		s.TimeOver <- false
	case <-time.After(15 * time.Second):
		if s.Content == Bomb {
			player.Score -= 5
		} else if s.Content == Apple {
			player.Score += 10
		}
		s.Content = Empty
		s.TimeOver <- true
	}
}

func CreateBoard(width int, height int) []*Cell {
	board := make([]*Cell, width*height)
	for i := range board {
		board[i] = newCell()
	}
	return board
}
