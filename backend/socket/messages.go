package socket

import "gamefr/game"

type receiveInitMessage struct {
	PlayerRole int `json:"mode"`
}

type receiveMsgPlayer1 struct {
	Pos           int `json:"pos"`
	ButtonClicked int `json:"button_clicked"`
}

type receiveMsgPlayer2 struct {
	Pos     int `json:"pos"`
	Content int `json:"contains"`
}

type response struct {
	RoomState   int `json:"room_state"`
	PlayerScore int `json:"score"`
	ExtraData   any `json:"data"`
}

func CreateResponse(room *Room, player *game.Player, data any) response {
	return response{
		RoomState:   room.Status,
		PlayerScore: player.Score,
		ExtraData:   data,
	}
}

type sendOnBombAppleButtonClicked struct {
	BombCount  int `json:"bomb_count"`
	AppleCount int `json:"apple_count"`
}

type sendOnActiveCellsButtonClicked struct {
	ActiveCells []int `json:"active_cells"`
}

type sendPlayer2 struct {
	Pos     int `json:"pos"`
	Content int `json:"content"`
}

type gameOverMsg struct {
	Result string `json:"result"`
	Note   string `json:"note"`
}

func createGameOverMsg(room *Room, player *game.Player, result string, note string) response {
	data := gameOverMsg{
		Result: result,
		Note:   note,
	}
	return CreateResponse(room, player, data)
}
