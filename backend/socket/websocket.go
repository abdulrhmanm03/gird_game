package socket

import (
	"errors"
	"gamefr/game"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for simplicity
	},
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
			err = handleMsgFromPlayer1(room, conn)
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

func getRoomId(min int, max int) int {
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	return r.Intn(max) + min
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

	roomId := getRoomId(1, 1000000)
	room, err := findOrCreateRoom(&player, roomId)
	if err != nil {
		log.Println("failed to creat a room")
		return nil, nil, err
	}

	return &player, room, nil
}

type InitMessageToSend struct {
	RoomState int `json:"room_state"`
}

func sendInitMsg(room *Room) error {
	msgToSend := InitMessageToSend{RoomState: room.Status}
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

type MsgFromPlayer1 struct {
	Pos int `json:"pos"`
}

type Player1ToPlayer1Msg struct {
	RoomState     int `json:"room_state"`
	SquereIdx     int `json:"squere_index"`
	SquereContent int `json:"squere_content"`
	Player1Score  int `json:"score"`
}

type Player1ToPlayer2Msg struct {
	RoomState    int `json:"room_state"`
	Pos          int `json:"pos"`
	Player2Score int `json:"score"`
}

// TODO: refactor this
func handleMsgFromPlayer1(room *Room, conn *websocket.Conn) error {

	// handle the message sent from player 1
	var msgFromPlayer1 MsgFromPlayer1
	err := conn.ReadJSON(&msgFromPlayer1)
	if err != nil {
		log.Println("read:", err)
		return err
	}
	log.Println("player1: ", msgFromPlayer1.Pos) // logging

	// game logic
	squere := room.Board[msgFromPlayer1.Pos]
	squereContent := squere.Content

	if squereContent == 1 {
		room.Player1.Score -= 5
		squere.Clicked <- true
	} else if squereContent == 2 {
		room.Player1.Score += 5
		squere.Clicked <- true
	}

	// send message to player 1
	player1ToPlayer1Msg := Player1ToPlayer1Msg{
		RoomState:     room.Status,
		SquereIdx:     msgFromPlayer1.Pos,
		SquereContent: squereContent,
		Player1Score:  room.Player1.Score,
	}
	err = room.Player1.Conn.WriteJSON(player1ToPlayer1Msg)
	if err != nil {
		log.Println("write:", err)
	}

	// send message to player 2
	player1ToPlayer2Msg := Player1ToPlayer2Msg{
		RoomState:    room.Status,
		Pos:          msgFromPlayer1.Pos,
		Player2Score: room.Player2.Score,
	}
	if room.Player2 != nil {
		err = room.Player2.Conn.WriteJSON(player1ToPlayer2Msg)
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
	squere := room.Board[msg.Pos]
	if squere.Content == 0 {
		squere.Content = msg.Contains
		go squere.Run(room.Player2, msg.Pos)
	}
	return nil
}
