package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

var grid [25]bool
var mu sync.Mutex

// Define a WebSocket upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for simplicity
	},
}

// Clients slice to hold connected clients
var clients = make(map[*websocket.Conn]bool)

// Handle WebSocket connections
func handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error during connection upgrade:", err)
		return
	}
	defer conn.Close()

	println(conn)
	clients[conn] = true
	// Send the current grid to the new client
	if err := conn.WriteJSON(grid); err != nil {
		fmt.Println("Error sending initial grid:", err)
		return
	}

	for {
		var newGrid [25]bool
		if err := conn.ReadJSON(&newGrid); err != nil {
			fmt.Println("Error reading from client:", err)
			break
		}

		// Update shared grid and broadcast it to all clients
		mu.Lock()
		grid = newGrid
		mu.Unlock()

		for client := range clients {
			if err := client.WriteJSON(grid); err != nil {
				fmt.Println("Error sending grid to client:", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func main() {
	http.HandleFunc("/ws", handleConnection)
	fmt.Println("WebSocket server is running on ws://localhost:8080/ws")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("Error starting server: " + err.Error())
	}
}
