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
	roomForMode1 *Room
	roomForMode2 *Room
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
		roomForMode2 = room

		return room, nil
	} else if player.Role == 2 {
		room.Player1 = nil
		room.Player2 = player
		roomForMode1 = room

		return room, nil
	}
	return nil, errors.New("not valid player role")
}

func addPlayerToRoom(player *game.Player, room *Room) (*Room, error) {
	if room.Status == 0 {
		return nil, errors.New("Room is full")
	}
	if room.Player1 == nil {
		room.Player1 = player
		room.Status = 0
		return room, nil
	}
	if room.Player2 == nil {
		room.Player2 = player
		room.Status = 0
		return room, nil
	}
	return nil, errors.New("something wrong happend")
}

func FindOrCreateRoom(player *game.Player, roomId int) (*Room, error) {
	if player.Role == 1 && roomForMode1 != nil && roomForMode1.Status == 1 {
		room, err := addPlayerToRoom(player, roomForMode1)
		if err != nil {
			return nil, err
		}
		return room, nil
	}
	if player.Role == 2 && roomForMode2 != nil && roomForMode2.Status == 1 {
		room, err := addPlayerToRoom(player, roomForMode2)
		if err != nil {
			return nil, err
		}
		return room, nil
	}
	room, err := createRoom(player, roomId)
	if err != nil {
		return nil, err
	}
	if player.Role == 1 {
		roomForMode2 = room
	} else {
		roomForMode2 = room
	}

	return room, nil
}
