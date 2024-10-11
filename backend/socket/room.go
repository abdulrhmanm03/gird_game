package socket

import (
	"errors"
	"gamefr/game"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Room struct {
	Id      int
	Player1 *game.Player
	Player2 *game.Player
	Status  int
	Board   []*game.Squere
}

// enum for room state
const (
	gameOver = iota
	waiting
	active
)

const roomTimeInMinutes = 3

var (
	roomsForMode1 = make(map[int]*Room)
	roomsForMode2 = make(map[int]*Room)
	activeRooms   = make(map[int]*Room)
)

var (
	width  = 5
	height = 5
)

func createRoom(player *game.Player, id int) (*Room, error) {
	room := &Room{
		Id:     id,
		Status: waiting,
		Board:  game.CreateBoard(width, height),
	}
	if player.Role == 1 {
		room.Player1 = player
		room.Player2 = nil
		roomsForMode2[room.Id] = room

		return room, nil
	} else if player.Role == 2 {
		room.Player1 = nil
		room.Player2 = player
		roomsForMode1[room.Id] = room

		return room, nil
	}
	return nil, errors.New("not valid player role")
}

func addPlayerToRoom(player *game.Player, room *Room) (*Room, error) {
	if room.Status == active {
		return nil, errors.New("Room is full")
	}
	if room.Player1 == nil {
		// test the other player connection if not connected delete the room
		err := room.Player2.Conn.WriteMessage(websocket.PingMessage, nil)
		if err != nil {
			delete(roomsForMode2, room.Id)
			return nil, err
		}
		room.Player1 = player
		room.Status = active
		activeRooms[room.Id] = room
		go monitorRoomConnection(room)
		go startRoomTimer(roomTimeInMinutes, room)
		return room, nil
	}
	if room.Player2 == nil {
		// test the other player connection if not connected delete the room
		err := room.Player1.Conn.WriteMessage(websocket.PingMessage, nil)
		if err != nil {
			delete(roomsForMode1, room.Id)
			return nil, err
		}
		room.Player2 = player
		room.Status = active
		activeRooms[room.Id] = room
		go monitorRoomConnection(room)
		go startRoomTimer(roomTimeInMinutes, room)
		return room, nil
	}
	return nil, errors.New("something wrong happend")
}

func findOrCreateRoom(player *game.Player, roomId int) (*Room, error) {
	var rooms map[int]*Room
	if player.Role == 1 {
		rooms = roomsForMode1
	} else {
		rooms = roomsForMode2
	}

	for _, room := range rooms {
		if room.Status == waiting {
			room, err := addPlayerToRoom(player, room)
			if err != nil {
				continue
			}
			return room, nil
		}
	}

	room, err := createRoom(player, roomId)
	if err != nil {
		return nil, err
	}

	return room, nil
}

type gameOverMsg struct {
	Room_state int    `json:"room_state"`
	Result     string `json:"result"`
	Note       string `json:"note"`
}

func createGameOverMsg(result string, note string) gameOverMsg {
	return gameOverMsg{
		Room_state: gameOver,
		Result:     result,
		Note:       note,
	}
}

func monitorRoomConnection(room *Room) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if room.Status == active {
			err1 := room.Player1.Conn.WriteMessage(websocket.PingMessage, nil)
			err2 := room.Player2.Conn.WriteMessage(websocket.PingMessage, nil)

			note := "The other player disconnected"
			result := "win"
			msg := createGameOverMsg(result, note)
			if err1 != nil {
				_ = room.Player2.Conn.WriteJSON(msg)
			}

			if err2 != nil {
				_ = room.Player1.Conn.WriteJSON(msg)
			}

			if err1 != nil || err2 != nil {
				delete(activeRooms, room.Id)
				return
			}
		}
	}
}

func writeResultsToPlayers(winner *game.Player, loser *game.Player, draw bool) {
	note := "Game over"

	if draw {
		tieMsg := createGameOverMsg("draw", note)
		_ = winner.Conn.WriteJSON(tieMsg)
		_ = loser.Conn.WriteJSON(tieMsg)
		return
	}

	winnerMsg := createGameOverMsg("win", note)
	loserMsg := createGameOverMsg("lose", note)

	_ = winner.Conn.WriteJSON(winnerMsg)
	_ = loser.Conn.WriteJSON(loserMsg)
}

func startRoomTimer(minutes int, room *Room) {
	player1 := room.Player1
	player2 := room.Player2

	defer player1.Conn.Close()
	defer player2.Conn.Close()

	log.Println(room.Id, " started")
	<-time.After(time.Duration(minutes) * time.Minute)
	log.Println(room.Id, " ended")

	room.Status = gameOver

	if player1.Score > player2.Score {
		writeResultsToPlayers(room.Player1, room.Player2, false)
	} else if player2.Score > player1.Score {
		writeResultsToPlayers(player2, player1, false)
	} else {
		writeResultsToPlayers(player1, player2, true)
	}
}
