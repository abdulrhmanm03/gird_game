package main

import (
	"gamefr/socket"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/ws", socket.WebsocketHandler)

	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Println("failed to start server:", err)
		return
	}

	log.Println("server started")
}
