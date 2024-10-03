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
	Board   [25]int
}

var (
	roomForMode1     = make(map[int]*Room)
	roomForMode1List = []int{}
	roomForMode2     = make(map[int]*Room)
	roomForMode2List = []int{}
)

func createRoom(player *game.Player, id int) (*Room, error) {
	room := &Room{
		Id:     id,
		Status: 1,
		Board:  [25]int{},
	}
	if player.Role == 1 {
		room.Player1 = player
		room.Player2 = nil
		roomForMode2[room.Id] = room
		roomForMode2List = append(roomForMode2List, room.Id)

		return room, nil
	} else if player.Role == 2 {
		room.Player1 = nil
		room.Player2 = player
		roomForMode1[room.Id] = room
		roomForMode1List = append(roomForMode1List, room.Id)

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
			delete(roomForMode2, room.Id)
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
			delete(roomForMode1, room.Id)
			return nil, err
		}
		room.Player2 = player
		room.Status = 0
		return room, nil
	}
	return nil, errors.New("something wrong happend")
}

func findOrCreateRoom(player *game.Player, roomId int) (*Room, error) {
	var roomList []int
	var rooms map[int]*Room
	if player.Role == 1 {
		roomList = roomForMode1List
		rooms = roomForMode1
	} else {
		roomList = roomForMode2List
		rooms = roomForMode2
	}

	for _, id := range roomList {
		if rooms[id].Status == 1 {
			room, err := addPlayerToRoom(player, rooms[id])
			if err != nil {
				return nil, err
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
