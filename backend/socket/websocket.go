package socket

import (
	"gamefr/game"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var mu sync.Mutex

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

	err = gameLoop(player, room)
	if err != nil {
		return
	}

}

func gameLoop(player *game.Player, room *Room) error {
	for {
		switch player.Role {
		case 1:
			err := handlePlayer1(room, player.Conn)
			if err != nil {
				return err
			}
		case 2:
			err := handlePlayer2(room, player.Conn)
			if err != nil {
				return err
			}
		}
	}
}

func generateRoomId(min int, max int) int {
	mu.Lock()
	defer mu.Unlock()
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	return r.Intn(max) + min
}

func handleInitMsg(conn *websocket.Conn) (*game.Player, *Room, error) {
	var initMsg receiveInitMessage
	err := conn.ReadJSON(&initMsg)
	if err != nil {
		log.Println("Error while reading message:", err)
		return nil, nil, err
	}

	player := game.CreatePlayer(initMsg.PlayerRole, conn)

	roomId := generateRoomId(1, 1000000)
	room, err := findOrCreateRoom(&player, roomId)
	if err != nil {
		log.Println("failed to creat a room")
		return nil, nil, err
	}

	return &player, room, nil
}

func sendInitMsg(room *Room) error {
	res := CreateResponse(room, room.Player1, nil)
	log.Println(res)

	if room.Player1 != nil {
		err := room.Player1.Conn.WriteJSON(res)
		if err != nil {
			log.Println("write:", err)
			return err
		}
	}
	if room.Player2 != nil {
		err := room.Player2.Conn.WriteJSON(res)
		if err != nil {
			log.Println("write:", err)
			return err
		}
	}
	return nil
}
