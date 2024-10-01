package main

import (
	"encoding/json"
	"gamefr/game"
	"gamefr/websocket"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for simplicity
	},
}

type InitMessage struct {
	RoomID int `json:"room_id"`
	Mode   int `json:"mode"`
}

type Mode1Message struct {
	Pos int `json:"pos"`
}

type Mode2Message struct {
	Pos      int `json:"pos"`
	Contains int `json:"contains"`
}

type MessageToSend struct {
	Room_state int     `json:"room_state"`
	Board      [25]int `json:"board"`
}

func websocketHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Failed to upgrade to WebSocket:", err)
		return
	}
	defer conn.Close()

	mt, message, err := conn.ReadMessage()
	if err != nil {
		log.Println("Error while reading message:", err)
		return
	}

	var msg InitMessage
	err = json.Unmarshal(message, &msg)
	if err != nil {
		log.Println("error while convirting to json:", err)
		return
	}

	log.Printf("recv: %d", msg.RoomID)
	player := game.CreatePlayer(msg.Mode, conn)
	room, err := socket.FindOrCreateRoom(&player, msg.RoomID)
	if err != nil {
		log.Println("failed to creat a room")
		return
	}

	msgToSend := MessageToSend{Room_state: room.Status, Board: room.Board}
	MsgJson, err := json.Marshal(msgToSend)
	if err != nil {
		log.Println("cant convert to json")
		return
	}
	if room.Player1 != nil {

		err = room.Player1.Conn.WriteMessage(mt, MsgJson)
		if err != nil {
			log.Println("write:", err)
		}
	}
	if room.Player2 != nil {

		err = room.Player2.Conn.WriteMessage(mt, MsgJson)
		if err != nil {
			log.Println("write:", err)
		}
	}
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		if player.Role == 1 {
			var msg Mode1Message
			err = json.Unmarshal(message, &msg)
			if err != nil {
				log.Println("error while convirting to json:", err)
				return
			}
			log.Println("player1: ", msg.Pos)
			log.Println(room.Board)
			room.Board[msg.Pos] = 0
			msgToSend := MessageToSend{Room_state: 1, Board: room.Board}
			MsgJson, err := json.Marshal(msgToSend)
			if err != nil {
				log.Println("cant convert to json")
				return
			}
			if room.Player2 != nil {
				err = room.Player2.Conn.WriteMessage(mt, MsgJson)
				if err != nil {
					log.Println("write:", err)
				}
			} else {
				log.Println("no player 2")
			}
		}
		if player.Role == 2 {
			var msg Mode2Message
			err = json.Unmarshal(message, &msg)
			if err != nil {
				log.Println("error while convirting to json:", err)
				return
			}
			log.Println("player2: ", msg.Pos, msg.Contains)
			room.Board[msg.Pos] = msg.Contains
			log.Println(room.Board)
		}
	}
}

type RoomIdRequest struct {
	Mode int `json:"mode" binding:"required"`
}

func getRoomId(c *gin.Context) {
	var req RoomIdRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"room_id": 9999}) // room id hard coded for now
}

func main() {
	r := gin.Default()
	corsConfig := cors.Config{
		AllowOrigins:     []string{"*"}, // Allow all origins
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}

	// Apply the CORS middleware
	r.Use(cors.New(corsConfig))
	r.POST("/getRoomId", getRoomId)
	r.GET("/ws", websocketHandler)
	err := r.Run(":3000")
	if err != nil {
		log.Println("failed to start server")
	}
	log.Println("server started")
}
