package main

import (
	"gamefr/socket"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/ws", socket.WebsocketHandler)

	log.Println("server started")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Println("failed to start server:", err)
		return
	}
}
