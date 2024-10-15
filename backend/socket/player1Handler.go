package socket

import (
	"errors"
	"gamefr/game"
	"log"

	"github.com/gorilla/websocket"
)

const (
	cellClicked = iota
	bombAppleCountButton
	activeCellButton
)

func sendMsgToplayer1(room *Room, pos int, cellContent int) error {

	data := sendPlayer2{
		Pos:     pos,
		Content: cellContent,
	}

	res := CreateResponse(room, room.Player1, data)

	err := room.Player1.Conn.WriteJSON(res)
	if err != nil {
		log.Println("write:", err)
		return err
	}
	return nil
}

func sendMsgToplayer2(room *Room, pos int) error {

	data := sendPlayer2{
		Pos:     pos,
		Content: game.Empty,
	}

	res := CreateResponse(room, room.Player2, data)

	if room.Player2 != nil {
		err := room.Player2.Conn.WriteJSON(res)
		if err != nil {
			log.Println("write:", err)
			return err
		}
	} else {
		log.Println("no player 2")
		return errors.New("player 2 disconnected")
	}
	return nil
}

func handlePlayer1CellClick(room *Room, pos int) error {

	cell := room.Board[pos]
	cellContent := cell.Content

	if cellContent == game.Bomb {
		room.Player1.Score -= 10
		cell.Clicked <- true
	} else if cellContent == game.Apple {
		room.Player1.Score += 5
		cell.Clicked <- true
	}
	cell.Content = game.Empty

	err := sendMsgToplayer1(room, pos, cellContent)
	if err != nil {
		return err
	}
	err = sendMsgToplayer2(room, pos)
	if err != nil {
		return err
	}
	return nil
}

func handleButtonClicked(room *Room, buttonClicked int) error {
	room.Player1.Score -= 5

	var err error
	if buttonClicked == bombAppleCountButton {
		err = onBombAppleButtonClick(room)
	} else if buttonClicked == activeCellButton {
		err = onActiveCellsButtonClicke(room)
	} else {
		return errors.New("not a valid button")
	}
	return err
}

func onBombAppleButtonClick(room *Room) error {

	bombCount := 0
	appleCount := 0
	for _, v := range room.Board {
		if v.Content == game.Bomb {
			bombCount++
		} else if v.Content == game.Apple {
			appleCount++
		}
	}

	data := sendOnBombAppleButtonClicked{
		BombCount:  bombCount,
		AppleCount: appleCount,
	}

	res := CreateResponse(room, room.Player1, data)

	err := room.Player1.Conn.WriteJSON(res)
	if err != nil {
		log.Println("write: faild to send msg to player 1")
		return err
	}

	return nil
}

func onActiveCellsButtonClicke(room *Room) error {

	activeCells := make([]int, 0, 25)
	for i, v := range room.Board {
		if v.Content == 1 || v.Content == 2 {
			activeCells = append(activeCells, i)
		}
	}

	data := sendOnActiveCellsButtonClicked{
		ActiveCells: activeCells,
	}

	res := CreateResponse(room, room.Player1, data)

	err := room.Player1.Conn.WriteJSON(res)
	if err != nil {
		log.Println("write: faild to send msg to player 1")
		return err
	}
	return nil
}

func handlePlayer1(room *Room, conn *websocket.Conn) error {

	var msgFromPlayer1 receiveMsgPlayer1
	err := conn.ReadJSON(&msgFromPlayer1)
	if err != nil {
		log.Println("read:", err)
		return err
	}

	if msgFromPlayer1.ButtonClicked == cellClicked {
		err := handlePlayer1CellClick(room, msgFromPlayer1.Pos)
		if err != nil {
			return err
		}
	} else {
		err := handleButtonClicked(room, msgFromPlayer1.ButtonClicked)
		if err != nil {
			return err
		}
	}

	return nil
}
