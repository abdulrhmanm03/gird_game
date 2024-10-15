package socket

import (
	"errors"
	"gamefr/game"
	"log"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
)

type Room struct {
	Id      int
	Player1 *game.Player
	Player2 *game.Player
	Status  int
	Board   []*game.Cell
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
	switch player.Role {
	case 1:
		room.Player1 = player
		roomsForMode2[room.Id] = room

		return room, nil
	case 2:
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
		initiateRoom(room)
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
		initiateRoom(room)
		return room, nil
	}
	return nil, errors.New("something wrong happend with adding player to the room")
}

func findOrCreateRoom(player *game.Player, roomId int) (*Room, error) {
	var availableRooms map[int]*Room
	if player.Role == 1 {
		availableRooms = roomsForMode1
	} else {
		availableRooms = roomsForMode2
	}

	for _, room := range availableRooms {
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

func monitorRoomConnection(room *Room) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if room.Status == active {
			mu.Lock()
			err1 := room.Player1.Conn.WriteMessage(websocket.PingMessage, nil)
			err2 := room.Player2.Conn.WriteMessage(websocket.PingMessage, nil)

			note := "The other player disconnected"
			result := "win"
			if err1 != nil {
				room.Status = gameOver
				msg := createGameOverMsg(room, room.Player2, result, note)
				log.Println(note)
				_ = room.Player2.Conn.WriteJSON(msg)
			}

			if err2 != nil {
				room.Status = gameOver
				msg := createGameOverMsg(room, room.Player1, result, note)
				log.Println(note)
				_ = room.Player1.Conn.WriteJSON(msg)
			}

			if err1 != nil || err2 != nil {
				delete(activeRooms, room.Id)
				mu.Unlock()
				return
			}
			mu.Unlock()
		}
	}
}

func writeResultsToPlayers(room *Room, winner *game.Player, loser *game.Player, draw bool) {
	mu.Lock()
	note := "Game over"

	defer winner.Conn.Close()
	defer loser.Conn.Close()
	defer mu.Unlock()

	if draw {
		tieMsg := createGameOverMsg(room, winner, "draw", note)
		_ = winner.Conn.WriteJSON(tieMsg)
		_ = loser.Conn.WriteJSON(tieMsg)
		return
	}

	winnerMsg := createGameOverMsg(room, winner, "win", note)
	loserMsg := createGameOverMsg(room, loser, "lose", note)

	_ = winner.Conn.WriteJSON(winnerMsg)
	_ = loser.Conn.WriteJSON(loserMsg)
}

func startRoomTimer(minutes int, room *Room) {
	log.Println(room.Id, " room started")
	player1 := room.Player1
	player2 := room.Player2

	defer player1.Conn.Close()
	defer player2.Conn.Close()
	defer log.Println(room.Id, " room finished")

	<-time.After(time.Duration(minutes) * time.Minute)

	room.Status = gameOver

	if player1.Score > player2.Score {
		writeResultsToPlayers(room, room.Player1, room.Player2, false)
	} else if player2.Score > player1.Score {
		writeResultsToPlayers(room, player2, player1, false)
	} else {
		writeResultsToPlayers(room, player1, player2, true)
	}
}

func plantRandomCells(room *Room) {
	ticker := time.NewTicker(10 * time.Second)
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	for range ticker.C {
		cellPos := r.Intn(len(room.Board))
		appleOrBomb := r.Intn(2) + 1

		cell := room.Board[cellPos]

		if cell.Content == 0 {

			cell.Content = appleOrBomb
			go cell.Run(room.Player2, cellPos)
			go sendWhenSquereTimeEnd(room, cell, cellPos)

			mu.Lock()
			data := sendPlayer2{
				Pos:     cellPos,
				Content: cell.Content,
			}

			res := CreateResponse(room, room.Player2, data)
			if room.Player2 != nil {
				err := room.Player2.Conn.WriteJSON(res)
				if err != nil {
					log.Println("write:", err)
					mu.Unlock()
					return
				}
			}
			mu.Unlock()
		}
	}
}

func initiateRoom(room *Room) {
	room.Status = active
	activeRooms[room.Id] = room
	go monitorRoomConnection(room)
	go startRoomTimer(roomTimeInMinutes, room)
	go plantRandomCells(room)
}
