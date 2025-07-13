package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	//shared resources
	clients = make(map[*websocket.Conn]bool)
	//Mutex to protect the shared resource
	clientsMu sync.RWMutex
	upgrader  = websocket.Upgrader{}
)

func main() {
	http.HandleFunc("/ws", handleConnections)

	go broadcastLoop()

	fmt.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Accepts new websocket clients
func handleConnections(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true } //Allow all origins
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer ws.Close()

	//Add client to shared map while protecting it with mutex
	clientsMu.Lock()
	clients[ws] = true
	clientsMu.Unlock()

	log.Println("New client connected")

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println("client disconnected")
			break
		}

		//Send to broadcast goroutine
		messages <- msg
	}

	//Remove client when disconnected
	clientsMu.Lock()
	delete(clients, ws)
	clientsMu.Unlock()
}

var messages = make(chan []byte)

// Reads messages from the `messages` channel and broadcasts them
// Reads messages from the `messages` channel and broadcasts them
func broadcastLoop() {
	for msg := range messages {
		// Lock for reading
		clientsMu.RLock()
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Println("Write error:", err)
				client.Close()

				// Remove dead client
				clientsMu.RUnlock()
				clientsMu.Lock()
				delete(clients, client)
				clientsMu.Unlock()
				clientsMu.RLock()
			}
		}
		clientsMu.RUnlock()
	}
}
