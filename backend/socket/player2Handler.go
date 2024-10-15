package socket

import (
	"gamefr/game"
	"log"

	"github.com/gorilla/websocket"
)

func handlePlayer2(room *Room, conn *websocket.Conn) error {
	var msg receiveMsgPlayer2
	err := conn.ReadJSON(&msg)
	if err != nil {
		log.Println("read:", err)
		return err
	}
	log.Println("player2: ", msg.Pos, msg.Content)
	cell := room.Board[msg.Pos]
	if cell.Content == 0 {
		cell.Content = msg.Content
		go cell.Run(room.Player2, msg.Pos)
		go sendWhenSquereTimeEnd(room, cell, msg.Pos)
		room.Player2.Score -= 5
	}
	return nil
}
func sendWhenSquereTimeEnd(room *Room, cell *game.Cell, pos int) {
	timeOver := <-cell.TimeOver
	if timeOver {

		data := sendPlayer2{
			Pos:     pos,
			Content: cell.Content,
		}
		mu.Lock()
		res := CreateResponse(room, room.Player2, data)
		err := room.Player2.Conn.WriteJSON(res)
		if err != nil {
			log.Println("write faild to write to player")
		}
		mu.Unlock()

	}
}
