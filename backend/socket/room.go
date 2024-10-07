package socket

import (
	"errors"
	"gamefr/game"
)

type Room struct {
	Id      int
	Player1 *game.Player
	Player2 *game.Player
	Status  int
	Board   []*game.Squere
}

var (
	roomsForMode1 = make(map[int]*Room)
	roomsForMode2 = make(map[int]*Room)
)

var (
	width  = 5
	height = 5
)

func createRoom(player *game.Player, id int) (*Room, error) {
	room := &Room{
		Id:     id,
		Status: 1,
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

type testConn struct {
	Test string `json:"test"`
}

func addPlayerToRoom(player *game.Player, room *Room) (*Room, error) {
	if room.Status == 0 {
		return nil, errors.New("Room is full")
	}
	if room.Player1 == nil {
		// test the other player connection if not connected delete the room
		err := room.Player2.Conn.WriteJSON(testConn{Test: "test"})
		if err != nil {
			delete(roomsForMode2, room.Id)
			return nil, err
		}
		room.Player1 = player
		room.Status = 0
		return room, nil
	}
	if room.Player2 == nil {
		// test the other player connection if not connected delete the room
		err := room.Player1.Conn.WriteJSON(testConn{Test: "test"})
		if err != nil {
			delete(roomsForMode1, room.Id)
			return nil, err
		}
		room.Player2 = player
		room.Status = 0
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
		if room.Status == 1 {
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
