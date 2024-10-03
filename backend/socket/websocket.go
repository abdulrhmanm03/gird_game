package socket

import (
	"errors"
	"gamefr/game"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for simplicity
	},
}

type InitMessage struct {
	Mode int `json:"mode"`
}

func handleInitMsg(conn *websocket.Conn) (*game.Player, *Room, error) {
	var initMsg InitMessage
	err := conn.ReadJSON(&initMsg)
	if err != nil {
		log.Println("Error while reading message:", err)
		return nil, nil, err
	}

	player := game.CreatePlayer(initMsg.Mode, conn)
	room, err := findOrCreateRoom(&player, 999)
	if err != nil {
		log.Println("failed to creat a room")
		return nil, nil, err
	}

	return &player, room, nil
}

type InitMessageToSend struct {
	Room_state int     `json:"room_state"`
	Board      [25]int `json:"board"`
}

func sendInitMsg(room *Room) error {
	msgToSend := InitMessageToSend{Room_state: room.Status}
	if room.Player1 != nil {
		err := room.Player1.Conn.WriteJSON(msgToSend)
		if err != nil {
			log.Println("write:", err)
			return err
		}
	}
	if room.Player2 != nil {
		err := room.Player2.Conn.WriteJSON(msgToSend)
		if err != nil {
			log.Println("write:", err)
			return err
		}
	}
	return nil
}

type Mode1Message struct {
	Pos int `json:"pos"`
}

type Mode1MessageToSend struct {
	Room_state int     `json:"room_state"`
	Board      [25]int `json:"board"`
}

func handlePlayer1(room *Room, conn *websocket.Conn) error {
	var msg Mode1Message
	err := conn.ReadJSON(&msg)
	if err != nil {
		log.Println("read:", err)
		return err
	}
	log.Println("player1: ", msg.Pos)
	log.Println(room.Board)
	room.Board[msg.Pos] = 0
	msgToSend := Mode1MessageToSend{Room_state: room.Status, Board: room.Board}
	if room.Player2 != nil {
		err = room.Player2.Conn.WriteJSON(msgToSend)
		if err != nil {
			log.Println("write:", err)
		}
	} else {
		log.Println("no player 2")
		return errors.New("player 2 left the room")
	}
	return nil
}

type Mode2Message struct {
	Pos      int `json:"pos"`
	Contains int `json:"contains"`
}

func handlePlayer2(room *Room, conn *websocket.Conn) error {
	var msg Mode2Message
	err := conn.ReadJSON(&msg)
	if err != nil {
		log.Println("read:", err)
		return err
	}
	log.Println("player2: ", msg.Pos, msg.Contains)
	room.Board[msg.Pos] = msg.Contains
	log.Println(room.Board)
	return nil
}

func WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not upgrade to websocket", http.StatusBadRequest)
		return
	}
	defer conn.Close()

	player, room, err := handleInitMsg(conn)
	if err != nil {
		return
	}

	err = sendInitMsg(room)
	if err != nil {
		return
	}

	for {
		if player.Role == 1 {
			err = handlePlayer1(room, conn)
			if err != nil {
				return
			}
		}
		if player.Role == 2 {
			err = handlePlayer2(room, conn)
			if err != nil {
				return
			}
		}
	}
}
